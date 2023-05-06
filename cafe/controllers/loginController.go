package controllers

import (
	"cafe/constants"
	"cafe/lib/database"
	"cafe/models"
	"database/sql"
	"main/config"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Login(c echo.Context) error {
	var credentials LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  constants.ErrorStatus,
			Message: err.Error(),
		})
	}

	// Connect to database
	config := config.LoadConfig()
	db, err := database.NewDB(config.DB_Username, config.DB_Password, config.DB_Host, config.DB_Port, config.DB_Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Failed to connect to database",
		})
	}
	defer db.Close()

	var user models.User
	row := db.Conn.QueryRow("SELECT * FROM users WHERE email = ? AND password = ?", credentials.Email, credentials.Password)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:  constants.ErrorStatus,
				Message: "Invalid email or password",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  constants.ErrorStatus,
				Message: "Failed to query database",
			})
		}
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(constants.JWTSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Failed to create token",
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Status:  constants.SuccessStatus,
		Message: "Login success",
		Token:   tokenString,
	})
}
