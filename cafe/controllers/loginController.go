package controllers

import (
	"cafe/config"
	"cafe/constants"
	"cafe/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Login(c echo.Context) error {
	var credentials models.LoginCredentials
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Invalid request body",
		})
	}
	cfg := config.LoadConfig()
	db, err := config.ConnectToDB(cfg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Failed to connect to database",
		})
	}
	defer func(db *gorm.DB) {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}(db)
	var user models.User
	err = db.Where("email = ? AND password = ?", credentials.Email, credentials.Password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	claims["userrole"] = user.Userrole
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(constants.JWTSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Failed to create token",
		})
	}
	return c.JSON(http.StatusOK, models.LoginResponse{
		Status:  constants.SuccessStatus,
		Message: "Login success",
		Token:   tokenString,
	})
}
