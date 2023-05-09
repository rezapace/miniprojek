package main

import (
	"cafe/config"
	"cafe/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitDB()

	// Initialize the Echo instance and pass the DB to the routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add the DB middleware
	e.Use(routes.DBMiddleware(config.DB))

	// Register the routes
	routes.RegisterRoutes(e)

	// Start the server
	e.Start(":8000")
}
