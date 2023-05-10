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

// TODO: membuat fungsi login
func Login(c echo.Context) error {
	// Deklarasi variabel credentials dengan tipe data models.LoginCredentials
	var credentials models.LoginCredentials
	// Jika terjadi error saat melakukan bind, maka return status bad request
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Invalid request body",
		})
	}

	// Gunakan koneksi database global
	db := config.DB

	// Deklarasi variabel user dengan tipe data models.User
	var user models.User
	// Query database berdasarkan email dan password yang dimasukkan
	err := db.Where("email = ? AND password = ?", credentials.Email, credentials.Password).First(&user).Error
	// Jika terjadi error saat query database
	if err != nil {
		// Jika error adalah record not found, maka return status bad request
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, models.Response{
				Status:  constants.ErrorStatus,
				Message: "Invalid email or password",
			})
		} else { // Jika error bukan record not found, maka return status internal server error
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  constants.ErrorStatus,
				Message: "Failed to query database",
			})
		}
	}
	// Buat token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	// Set claims
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["userrole"] = user.Userrole
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Signed string token
	tokenString, err := token.SignedString([]byte(constants.JWTSecret))
	// Jika terjadi error saat signed string token, maka return status internal server error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  constants.ErrorStatus,
			Message: "Failed to create token",
		})
	}
	// Return status OK dan token
	return c.JSON(http.StatusOK, models.LoginResponse{
		Status:  constants.SuccessStatus,
		Message: "Login success",
		Token:   tokenString,
	})
}
