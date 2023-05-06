package main

import (
	"cafe/config"
	"cafe/controllers"
	"cafe/lib/database"
	"cafe/middleware"
	"cafe/routes"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Set custom validator for request body validation
	e.Validator = &CustomValidator{validator: validator.New()}

	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize database connection
	db := database.InitDB(config.LoadConfig())

	// Set up routes
	routes.InitRoutes(e, controllers.NewOrderController(db), controllers.NewFoodController(db), controllers.NewUserController(db), controllers.NewLoginController(db))

	// Start server
	address := ":8000"
	fmt.Println("Server started at", address)
	e.Start(address)
}
