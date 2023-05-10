package routes

import (
	"cafe/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	// import jwt middleware
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
	e.GET("/users", controllers.GetUsers)        // data user
	e.GET("/users/:id", controllers.GetUserById) // data user by id

	// Routes for foods
	e.GET("/foods", controllers.GetFoods)          // data list makanan
	e.GET("/foods/:id", controllers.GetFoodById)   // data makanan by id
	e.POST("/foods", controllers.CreateFood)       // tambah makanan
	e.PUT("/foods/:id", controllers.UpdateFood)    // update makanan
	e.DELETE("/foods/:id", controllers.DeleteFood) // hapus makanan

	// Routes for orders
	e.GET("/orders", controllers.GetOrders)        // data list order
	e.GET("/orders/:id", controllers.GetOrderById) // data order by id
	e.POST("/orders", controllers.CreateOrder)     // tambah order
	e.PUT("/orders/:id", controllers.UpdateOrder)

	// Route for login dan registrasi
	e.POST("/login", controllers.Login)         //login
	e.POST("/register", controllers.CreateUser) // registrasi
}
