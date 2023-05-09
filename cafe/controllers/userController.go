package controllers

import (
	"main/config"
	"net/http"
	"strconv"

	"cafe/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetUsers returns all users
func GetUsers(c echo.Context) error {
	cfg := config.LoadConfig()
	db, err := config.ConnectToDB(cfg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
		})
	}
	defer func(db *gorm.DB) {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}(db)
	var users []models.User
	err = db.Select("id", "name", "email").Find(&users).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   users,
	})
}

// GetUserById returns a user based on ID
func GetUserById(c echo.Context) error {
	cfg := config.LoadConfig()
	db, err := config.ConnectToDB(cfg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
		})
	}
	defer func(db *gorm.DB) {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}(db)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
	}
	var user models.User
	err = db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  "error",
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   user,
	})
}

// CreateUser creates a new user
func CreateUser(c echo.Context) error {
	cfg := config.LoadConfig()
	db, err := config.ConnectToDB(cfg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to connect to database",
		})
	}
	defer func(db *gorm.DB) {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}(db)
	var form models.User
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}
	newUser := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
		Userrole: form.Userrole,
	}
	err = db.Create(&newUser).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to insert new user into database",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   newUser,
	})
}
