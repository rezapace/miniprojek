package controllers

import (
	"cafe/models"
	"fmt"
	"main/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GetOrders will retrieve all orders from database
func GetOrders(c echo.Context) error {
	// Get orders from database
	dbConfig := config.LoadConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DB_Username, dbConfig.DB_Password, dbConfig.DB_Host, dbConfig.DB_Port, dbConfig.DB_Name)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
		})
	}
	var orders []models.Order
	result := db.Preload("Details").Find(&orders)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   orders,
	})
}

// GetOrderById will retrieve order from database by given ID
func GetOrderById(c echo.Context) error {
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
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DB_Username, dbConfig.DB_Password, dbConfig.DB_Host, dbConfig.DB_Port, dbConfig.DB_Name)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
		})
	}
	var order models.Order
	result := db.Preload("Details").First(&order, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
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
	dbConfig := config.LoadConfig()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DB_Username, dbConfig.DB_Password, dbConfig.DB_Host, dbConfig.DB_Port, dbConfig.DB_Name)
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
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
	// Create new order
	order := &models.Order{
		UserID:     form.UserID,
		TotalPrice: total,
		Status:     form.Status,
	}
	err = models.CreateOrderInDB(db, order)
	if err != nil {
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
	err = models.CreateOrderDetailInDB(db, orderDetail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to create new order detail in database",
		})
	}
	// Return success response
	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: fmt.Sprintf("Order with ID %d has been created", order.ID),
		Data:    order,
	})
}
