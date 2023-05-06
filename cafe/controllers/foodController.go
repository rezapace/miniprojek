package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"cafe/models"

	"github.com/labstack/echo/v4"
)

type FoodController struct {
	DB *sql.DB
}

func (ctrl *FoodController) CreateFood(c echo.Context) error {
	var foodForm models.FoodForm
	if err := c.Bind(&foodForm); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&foodForm); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	food := models.Food{
		Name:        foodForm.Name,
		Description: foodForm.Description,
		Price:       foodForm.Price,
	}

	err := food.Insert(ctrl.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to insert food into database",
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Food has been created",
		Data:    food,
	})
}

func (ctrl *FoodController) GetFood(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid food ID",
		})
	}

	food, err := models.GetFoodByID(ctrl.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  "error",
				Message: "Food not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "error",
				Message: "Failed to query database",
			})
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   food,
	})
}

func (ctrl *FoodController) UpdateFood(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid food ID",
		})
	}

	var foodForm models.FoodForm
	if err := c.Bind(&foodForm); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(&foodForm); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	food := models.Food{
		ID:          id,
		Name:        foodForm.Name,
		Description: foodForm.Description,
		Price:       foodForm.Price,
	}

	err = food.Update(ctrl.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to update food in database",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Food has been updated",
		Data:    food,
	})
}

func (ctrl *FoodController) DeleteFood(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid food ID",
		})
	}

	err = models.DeleteFood(ctrl.DB, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to delete food from database",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Food has been deleted",
	})
}
