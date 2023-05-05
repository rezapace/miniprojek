package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"userrole"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterForm struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type Food struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type FoodForm struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=1"`
}

type FoodListResponse struct {
	Status string `json:"status"`
	Data   []Food `json:"data"`
}

type Order struct {
	ID         int           `json:"id"`
	UserID     int           `json:"user_id"`
	TotalPrice float64       `json:"total_price"`
	Status     string        `json:"status"`
	OrderTime  string        `json:"order_time"`
	Details    []OrderDetail `json:"details,omitempty"`
}

type OrderForm struct {
	UserID   int     `json:"user_id" validate:"required"`
	FoodID   int     `json:"food_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,min=1"`
	Status   string  `json:"status" validate:"required,oneof=proses selesai"`
}

type OrderDetail struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	FoodID   int     `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderStatusResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/cafe")
	if err != nil {
		fmt.Println("Error opening database:", err.Error())
		return
	}
	defer db.Close()

	// Initialize Echo instance
	e := echo.New()

	// Set up the custom validator
	e.Validator = &CustomValidator{validator: validator.New()}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Home page")
	})

	// User routes
	// ! REGISTER name,email,password,role
	e.POST("/register", func(c echo.Context) error {
		var form RegisterForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		// Check if user exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", form.Email).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "User with email already exists",
			})
		}
		// Insert new user
		result, err := db.Exec("INSERT INTO users (name, email, password, userrole) VALUES (?, ?, ?, ?)", form.Name, form.Email, form.Password, form.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to insert new user into database",
			})
		}
		id, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to retrieve user ID from database",
			})
		}
		// Return success response
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("User with ID %d has been created", id),
		})
	})

	// ! LOGIN email,password
	e.POST("/login", func(c echo.Context) error {
		var form LoginForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		// Check if user exists and password is correct
		var user User
		err := db.QueryRow("SELECT id, name, email, password, userrole FROM users WHERE email = ? AND password = ?", form.Email, form.Password).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusUnauthorized, Response{
					Status:  "error",
					Message: "Incorrect email or password",
				})
			} else {
				return c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Failed to query database",
				})
			}
		}
		// Return success response
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: "Login successful",
			Data:    user,
		})
	})

	// Food routes
	// ! GET daftar makanan
	e.GET("/foods", func(c echo.Context) error {
		rows, err := db.Query("SELECT id, name, description, price FROM food")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		defer rows.Close()

		foods := []Food{}
		for rows.Next() {
			var food Food
			if err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price); err != nil {
				return c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Failed to scan food rows",
				})
			}
			foods = append(foods, food)
		}

		return c.JSON(http.StatusOK, FoodListResponse{
			Status: "success",
			Data:   foods,
		})
	})

	// ! untuk menambahkan produk makanan baru
	e.POST("/foods", func(c echo.Context) error {
		var form FoodForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		result, err := db.Exec("INSERT INTO food (name, description, price) VALUES (?, ?, ?)", form.Name, form.Description, form.Price)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to insert new food into database",
			})
		}
		id, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to retrieve food ID from database",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Food with ID %d has been created", id),
		})
	})

	// Order routes
	// ! menampilkan detail pesanan berdasarkan ID pesanan
	e.GET("/orders/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid order ID",
			})
		}
		var order Order
		err = db.QueryRow("SELECT id, user_id, total_price, status, order_time FROM orders WHERE id = ?", id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.OrderTime)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, Response{
					Status:  "error",
					Message: "Order not found",
				})
			} else {
				return c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Failed to query database",
				})
			}
		}
		var details []OrderDetail
		rows, err := db.Query("SELECT id, order_id, food_id, quantity, price FROM order_details WHERE order_id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		defer rows.Close()
		for rows.Next() {
			var detail OrderDetail
			if err := rows.Scan(&detail.ID, &detail.OrderID, &detail.FoodID, &detail.Quantity, &detail.Price); err != nil {
				return c.JSON(http.StatusInternalServerError, Response{Status: "error",
					Message: "Failed to scan order detail rows",
				})
			}
			details = append(details, detail)
		}
		order.Details = details
		return c.JSON(http.StatusOK, Response{
			Status: "success",
			Data:   order,
		})
	})

	// ! untuk menambahkan pesanan baru (order)
	e.POST("/orders", func(c echo.Context) error {
		var form OrderForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		// Check if user exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", form.UserID).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		if count == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "User does not exist",
			})
		}
		// Check if food exists
		err = db.QueryRow("SELECT COUNT(*) FROM food WHERE id = ?", form.FoodID).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		if count == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Food does not exist",
			})
		}
		// Calculate total price
		total := form.Price * float64(form.Quantity)
		// Insert new order
		result, err := db.Exec("INSERT INTO orders (user_id, total_price, status) VALUES (?, ?, ?)", form.UserID, total, form.Status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to insert new order into database",
			})
		}
		// Get the last inserted order ID
		orderID, err := result.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to get the last inserted order ID",
			})
		}
		// Insert order detail
		_, err = db.Exec("INSERT INTO order_details (order_id, food_id, quantity, price) VALUES (?, ?, ?, ?)", orderID, form.FoodID, form.Quantity, form.Price)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to insert order detail into database",
			})
		}
		// Return success response
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Order with ID %d has been created", orderID),
		})
	})

	// !  fitur update order dalam sistem
	e.PUT("/orders/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid order ID",
			})
		}
		var form OrderForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{Status: "error",
				Message: err.Error(),
			})
		}
		// Check if order exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", id).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		if count == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Order does not exist",
			})
		}
		// Calculate total price
		total := form.Price * float64(form.Quantity)
		// Update order
		_, err = db.Exec("UPDATE orders SET user_id = ?, total_price = ?, status = ? WHERE id = ?", form.UserID, total, form.Status, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to update order in database",
			})
		}
		// Update order details
		_, err = db.Exec("UPDATE order_details SET food_id = ?, quantity = ?, price = ? WHERE order_id = ?", form.FoodID, form.Quantity, form.Price, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to update order details in database",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Order with ID %d has been updated", id),
		})
	})

	// ! fitur hapus order dalam sistem
	e.DELETE("/orders/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid order ID",
			})
		}
		// Check if order exists
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", id).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
		if count == 0 {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Order does not exist",
			})
		}
		// Delete order details
		_, err = db.Exec("DELETE FROM order_details WHERE order_id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to delete order details from database",
			})
		}
		// Delete order
		_, err = db.Exec("DELETE FROM orders WHERE id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to delete order from database",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Order with ID %d has been deleted", id),
		})
	})

	// ! mengubah data makanan berdasarkan ID makanan
	e.PUT("/food/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid food ID",
			})
		}
		var form FoodForm
		if err := c.Bind(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid request body",
			})
		}
		if err := c.Validate(&form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		_, err = db.Exec("UPDATE food SET name = ?, description = ?, price = ? WHERE id = ?", form.Name, form.Description, form.Price, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to update food in database",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Food with ID %d has been updated", id),
		})
	})

	// ! menghapus data makanan berdasarkan ID makanan
	e.DELETE("/food/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid food ID",
			})
		}
		_, err = db.Exec("DELETE FROM food WHERE id = ?", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Failed to delete food from database",
			})
		}
		return c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: fmt.Sprintf("Food with ID %d has been deleted", id),
		})
	})

	// ! untuk menampilkan status pesanan berdasarkan ID pesanan
	e.GET("/orders/status/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Invalid order ID",
			})
		}
		var orderStatus string
		err = db.QueryRow("SELECT status FROM orders WHERE id = ?", id).Scan(&orderStatus)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, Response{
					Status:  "error",
					Message: "Order not found",
				})
			} else {
				return c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Failed to query database",
				})
			}
		}
		return c.JSON(http.StatusOK, Response{
			Status: "success",
			Data: map[string]string{
				"status": orderStatus,
			},
		})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
