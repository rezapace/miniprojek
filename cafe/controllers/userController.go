package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"cafe/constants"
	"cafe/lib/database"
	"cafe/models"

	"github.com/labstack/echo/v4"
)

// GetUsers returns all users
func GetUsers(c echo.Context) error {
	// Get database connection
	db := database.DBInstance(constants.DB_Name)
	defer db.Close()

	// Query all users
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}
	defer rows.Close()

	// Create slice to store users
	var users []models.User

	// Loop through rows and append to slice
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  "error",
				Message: "Failed to scan rows",
			})
		}
		users = append(users, user)
	}

	// Return success response with users data
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   users,
	})
}

// GetUser returns a user based on ID
func GetUser(c echo.Context) error {
	// Get database connection
	db := database.DBInstance(constants.DBName)
	defer db.Close()

	// Get ID parameter from URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid user ID",
		})
	}

	// Query user with specified ID
	var user models.User
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
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

	// Return success response with user data
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   user,
	})
}

// CreateUser creates a new user
func CreateUser(c echo.Context) error {
	// Get database connection
	db := database.DBInstance(constants.DBName)
	defer db.Close()

	// Bind request body to UserForm struct
	var form models.UserForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	// Validate request body using UserFormValidator
	if err := c.Validate(&form); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	// Insert new user
	result, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", form.Name, form.Email, form.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to insert new user into database",
		})
	}

	// Get the last inserted user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to get the last inserted user ID",
		})
	}

	// Retrieve the newly created user from database
	var newUser models.User
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID).Scan(&newUser.ID, &newUser.Name, &newUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Failed to query database",
		})
	}

	// Return success response with the newly created user data
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   newUser,
	})
}
