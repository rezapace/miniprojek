package controllers

import (
	"cafe/lib/database"
	"cafe/models"
	"database/sql"
	"fmt"
	"main/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetOrders will retrieve all orders from database
func GetOrders(c echo.Context) error {
	// Get orders from database
	dbConfig := config.LoadConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.DB_Username, dbConfig.DB_Password, dbConfig.DB_Host, dbConfig.DB_Port, dbConfig.DB_Name)

	db, err := sql.Open("mysql", connStr)

	defer db.Close()

	rows, err := db.Query("SELECT id, user_id, total_price, status FROM orders")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}
	defer rows.Close()

	// Map query result to orders slice
	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "error",
				Message: "Failed to scan row",
			})
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   orders,
	})
}

// GetOrder will retrieve order from database by given ID
func GetOrder(c echo.Context) error {
	// Get order ID from request URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid order ID",
		})
	}

	// Get order from database
	dbConfig := config.LoadConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.DB_Username, dbConfig.DB_Password, dbConfig.DB_Host, dbConfig.DB_Port, dbConfig.DB_Name)

	db, err := sql.Open("mysql", connStr)

	defer db.Close()

	order := models.Order{}
	err = db.QueryRow("SELECT id, user_id, total_price, status FROM orders WHERE id = ?", id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  "error",
				Message: "Order not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}

	// Get order details from database
	orderDetails := make([]models.OrderDetail, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}

	order.Details = orderDetails
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   order,
	})
}

// CreateOrder function to create new order
func CreateOrder(c echo.Context) error {
	var form models.OrderForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}
	if err := c.Validate(&form); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	// Get database connection
	db := database.DBInstance(config.LoadConfig().DB_Name)
	defer db.Close()

	// Check if user exists
	user, err := models.GetUserById(db, form.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get user from database",
		})
	}
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "User does not exist",
		})
	}

	// Check if user exists
	user, err := models.GetUserById(db, form.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get user from database",
		})
	}
	if user == nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "User does not exist",
		})
	}

	// Check if food exists
	food, err := models.GetFoodById(db, form.FoodID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get food from database",
		})
	}
	if food == nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Food does not exist",
		})
	}

	// Calculate total price
	total := form.Price * float64(form.Quantity)

	// Begin database transaction
	tx, err := db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to begin database transaction",
		})
	}

	// Create new order
	order := &models.Order{
		UserID:     form.UserID,
		TotalPrice: total,
		Status:     form.Status,
	}
	err = models.CreateOrderInDB(tx, order)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to create new order in database",
		})
	}

	// Create new order detail
	orderDetail := &models.OrderDetail{
		OrderID:  order.ID,
		FoodID:   form.FoodID,
		Quantity: form.Quantity,
		Price:    form.Price,
	}
	err = models.CreateOrderDetailInDB(tx, orderDetail)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to create new order detail in database",
		})
	}

	// Commit database transaction
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to commit database transaction",
		})
	}

	// Return success response
	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: fmt.Sprintf("Order with ID %d has been created", order.ID),
		Data:    order,
	})
}
