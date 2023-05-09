package controllers

import (
	"cafe/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetFoods returns all food items
func GetFoods(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	foods, err := models.GetFood(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to retrieve foods from database",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   foods,
	})
}

// CreateFood creates a new food item
func CreateFood(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	food := new(models.Food)
	if err := c.Bind(food); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := models.CreateFoodInDB(db, food); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to create food in database",
		})
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"status": "success",
		"data":   food,
	})
}

// UpdateFood updates a food item by ID
func UpdateFood(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid food ID",
		})
	}
	food, err := models.GetFoodById(db, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to retrieve food from database",
		})
	}
	if food == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Food not found",
		})
	}
	if err := c.Bind(food); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := models.UpdateFoodInDB(db, food); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to update food in database",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   food,
	})
}

// DeleteFood deletes a food item by ID
func DeleteFood(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid food ID",
		})
	}
	food, err := models.GetFoodById(db, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to retrieve food from database",
		})
	}
	if food == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Food not found",
		})
	}
	if err := models.DeleteFoodInDB(db, id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to delete food in database",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": fmt.Sprintf("Food with ID %d has been deleted", id),
	})
}

// GetFoodById function to get food by ID
func GetFoodById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid food ID",
		})
	}
	db := c.Get("db").(*gorm.DB)                  // Get GORM db instance from context
	food, err := models.GetFoodById(db, uint(id)) // Convert id to uint
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get food from database",
		})
	}
	if food == nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "Food not found",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   food,
	})
}
