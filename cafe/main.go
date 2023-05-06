package main

import (
	"cafe/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", controllers.GetUsers)
	e.GET("/users/:id", controllers.GetUserById)
	e.POST("/users", controllers.CreateUser)

	// e.GET("/foods", controllers.GetFoods)
	// e.GET("/foods/:id", controllers.GetFood)
	// e.POST("/foods", controllers.CreateFood)
	// e.PUT("/foods/:id", controllers.UpdateFood)
	// e.DELETE("/foods/:id", controllers.DeleteFood)

	e.GET("/orders", controllers.GetOrders)
	e.GET("/orders/:id", controllers.GetOrder)
	e.POST("/orders", controllers.CreateOrder)

	e.POST("/login", controllers.Login)

	// Start server
	e.Start(":1323")
}
