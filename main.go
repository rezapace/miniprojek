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
	Role     string `json:"role"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RegisterForm struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
	Role     string `json:"role" form:"role" validate:"required"`
}

type Food struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Status      string `json:"status"`
}

// FoodForm struct for validation
type FoodForm struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int    `json:"price" validate:"required,min=1"`
}

// FoodListResponse struct
type FoodListResponse struct {
	Status string `json:"status"`
	Data   []Food `json:"data"`
}

type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	FoodID     int    `json:"food_id"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
}

// OrderForm struct for validation
type OrderForm struct {
	UserID   int    `json:"user_id" validate:"required"`
	FoodID   int    `json:"food_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Price    int    `json:"price" validate:"required,min=1"`
	Status   string `json:"status" validate:"required,oneof=pending done"`
}

// OrderListResponse struct
type OrderStatusResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	// Initialize database connection
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/food_ordering")
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

	// TODO: Implement the login, register, and other routes based on the new database schema

	// ! Route to show all user data
	e.GET("/users", func(c echo.Context) error {
		// Query all users from the database
		rows, err := db.Query("SELECT id, name, email, role FROM user")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}
		defer rows.Close()

		// Iterate through rows and append to users list
		users := []User{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
			}
			users = append(users, user)
		}

		// Check for any errors during row iteration
		if err := rows.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return users data in the response
		return c.JSON(http.StatusOK, users)
	})

	// ! Login route
	e.POST("/login", func(c echo.Context) error {
		// Bind and validate login form data
		form := new(LoginForm)
		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid login form data"})
		}
		if err := c.Validate(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid login form data"})
		}

		// Query user from database
		var user User
		row := db.QueryRow("SELECT id, name, email, password, role FROM user WHERE email = ?", form.Email)
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
			fmt.Printf("Error scanning row: %v\n", err) // Add this line for debugging
			return c.JSON(http.StatusUnauthorized, Response{"error", "Invalid email or password"})
		}
		if err == sql.ErrNoRows {
			fmt.Println("No user found with the given email")
			return c.JSON(http.StatusUnauthorized, Response{"error", "Invalid email or password"})
		}
		fmt.Printf("User from DB: %+v\n", user) // Add this line for debugging
		fmt.Printf("Form data: %+v\n", form)    // Add this line for debugging

		// Validate password
		if user.Password != form.Password {
			return c.JSON(http.StatusUnauthorized, Response{"error", "Invalid email or password"})
		}

		// Return success response
		return c.JSON(http.StatusOK, Response{"success", "Logged in successfully"})
	})

	// ! Register route
	e.POST("/register", func(c echo.Context) error {
		// Bind and validate register form data
		form := new(RegisterForm)
		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", fmt.Sprintf("Invalid register form data: %v", err)})
		}
		if err := c.Validate(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", fmt.Sprintf("Invalid register form data: %v", err)})
		}

		// Check if user already exists
		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ?", form.Email)
		if err := row.Scan(&count); err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, Response{"error", "Email already exists"})
		}

		// Insert user into the database
		_, err = db.Exec("INSERT INTO user (name, email, password, role) VALUES (?, ?, ?, ?)", form.Name, form.Email, form.Password, "user")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return success response
		return c.JSON(http.StatusOK, Response{"success", "User registered successfully"})
	})

	// ! Route for placing an order
	e.POST("/orders", func(c echo.Context) error {
		// Bind and validate order form data
		form := new(OrderForm)
		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid order form data"})
		}
		if err := c.Validate(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid order form data"})
		}

		// Calculate total price
		totalPrice := form.Quantity * form.Price

		// Insert order into database
		_, err = db.Exec("INSERT INTO `order` (user_id, food_id, quantity, price, total_price, status) VALUES (?, ?, ?, ?, ?, ?)", form.UserID, form.FoodID, form.Quantity, form.Price, totalPrice, form.Status)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return success response
		return c.JSON(http.StatusCreated, Response{"success", "Order placed successfully"})
	})

	// ! Route for adding a new menu item
	e.POST("/menu", func(c echo.Context) error {
		// Bind and validate food form data
		form := new(FoodForm)
		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid food form data"})
		}
		if err := c.Validate(form); err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid food form data"})
		}

		// Insert food item into database
		_, err = db.Exec("INSERT INTO `food` (name, description, price) VALUES (?, ?, ?)", form.Name, form.Description, form.Price)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return success response
		return c.JSON(http.StatusCreated, Response{"success", "Menu item added successfully"})
	})

	// ! Route for displaying menu items
	e.GET("/menu", func(c echo.Context) error {
		// Query all food items from the database
		rows, err := db.Query("SELECT id, name, description, price FROM `food`")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}
		defer rows.Close()

		// Iterate through rows and append to food list
		foodList := []Food{}
		for rows.Next() {
			var food Food
			err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
			}
			foodList = append(foodList, food)
		}

		// Check for any errors during row iteration
		if err := rows.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return response with food list
		return c.JSON(http.StatusOK, FoodListResponse{"success", foodList})
	})

	// OrderStatusResponse struct
	type OrderStatusResponse struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	// Route to check order status by ID
	e.GET("/order/:id/status", func(c echo.Context) error {
		// Get order ID from the URL parameter
		idParam := c.Param("id")
		orderID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid order ID"})
		}

		// Query order from the database by ID
		var order Order
		row := db.QueryRow("SELECT id, user_id, food_id, status FROM orders WHERE id = ?", orderID)
		if err := row.Scan(&order.ID, &order.UserID, &order.FoodID, &order.Status); err != nil {
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, Response{"error", "Order not found"})
			}
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Return order status in the response
		return c.JSON(http.StatusOK, OrderStatusResponse{"success", order.Status})
	})

	// Route to update order status to 'siap' by order ID for admin users
	e.PUT("/admin/order/:id/ready", func(c echo.Context) error {
		// Simulate getting the user role from the authenticated user (replace with actual authentication)
		userRole := "admin"

		// Ensure the user is an admin
		if userRole != "admin" {
			return c.JSON(http.StatusForbidden, Response{"error", "Access denied"})
		}

		// Get order ID from the URL parameter
		idParam := c.Param("id")
		orderID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{"error", "Invalid order ID"})
		}

		// Update the order status to 'siap' in the database
		result, err := db.Exec("UPDATE orders SET status = 'siap' WHERE id = ?", orderID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}

		// Check if any row was affected (order was found and updated)
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{"error", "Internal server error"})
		}
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound, Response{"error", "Order not found"})
		}

		// Return a success response
		return c.JSON(http.StatusOK, Response{"success", "Order status updated to 'siap'"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
