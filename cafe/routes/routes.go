package routes

import (
	"net/http"

	"cafe/config"
	"cafe/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Custom validator
	e.Validator = config.NewValidator()

	// Public routes
	e.POST("/login", controllers.Login)
	e.POST("/register", controllers.Register)

	// Restricted group
	r := e.Group("/api")

	// JWT middleware
	r.Use(middleware.JWT([]byte(config.JWT_SECRET)))

	// Food routes
	r.GET("/foods", controllers.GetFoods)
	r.GET("/foods/:id", controllers.GetFood)
	r.POST("/foods", controllers.CreateFood)
	r.PUT("/foods/:id", controllers.UpdateFood)
	r.DELETE("/foods/:id", controllers.DeleteFood)

	// Order routes
	r.GET("/orders", controllers.GetOrders)
	r.GET("/orders/:id", controllers.GetOrder)
	r.POST("/orders", controllers.CreateOrder)
	r.PUT("/orders/:id", controllers.UpdateOrderStatus)

	// User routes
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUser)

	// Authorization middleware for admin routes
	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(config.JWT_SECRET),
		Claims:      &controllers.JWTClaims{},
		TokenLookup: "header:" + echo.HeaderAuthorization,
	}))
	adminGroup.Use(controllers.AdminAuthorization)

	// Admin restricted routes
	adminGroup.GET("/orders", controllers.GetAdminOrders)
	adminGroup.GET("/users", controllers.GetAdminUsers)

	// Default error handling
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		response := controllers.Response{
			Status:  "error",
			Message: message,
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == "HEAD" {
				c.NoContent(code)
			} else {
				c.JSON(code, response)
			}
		}
	}

	return e
}
