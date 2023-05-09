package routes

import (
	"cafe/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func DBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

func RegisterRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes for users
	e.GET("/users", controllers.GetUsers)
	e.GET("/users/:id", controllers.GetUserById)
	e.POST("/users", controllers.CreateUser)

	// Routes for foods
	e.GET("/foods", controllers.GetFoods)
	e.GET("/foods/:id", controllers.GetFoodById)
	e.POST("/foods", controllers.CreateFood)
	e.PUT("/foods/:id", controllers.UpdateFood)
	e.DELETE("/foods/:id", controllers.DeleteFood)

	// Routes for orders
	e.GET("/orders", controllers.GetOrders)
	e.GET("/orders/:id", controllers.GetOrderById)
	e.POST("/orders", controllers.CreateOrder)

	// Route for login
	e.POST("/login", controllers.Login)
}
