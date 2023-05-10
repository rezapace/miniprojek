package controllers

import (
	"net/http"
	"strconv"

	"cafe/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TODO: mengambil semua data user
func GetUsers(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengambil semua data user dari database
	users, err := models.GetUsers(db)
	if err != nil {
		// Gagal mengambil data user dari database
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mengambil user dari database",
		})
	}
	// Berhasil mengambil data user dari database
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   users,
	})
}

// TODO: menampilkan user berdasarkan id
func GetUserById(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengubah string id menjadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Jika gagal mengubah string ke integer, maka return bad request
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "ID user tidak valid",
		})
	}

	// Mendapatkan data user berdasarkan id
	user, err := models.GetUserByID(db, uint(id))
	if err != nil {
		// Jika gagal mendapatkan data user, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mendapatkan user dari database",
		})
	}
	if user == nil {
		// Jika tidak ada data user yang ditemukan, maka return not found
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
	}
	// Return status OK dan data user
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   user,
	})
}

// TODO:  membuat user baru
func CreateUser(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Membuat variabel baru bertipe models.User
	user := new(models.User)
	// Mengecek apakah ada error saat binding data
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	// Mencoba membuat user baru di database
	if err := models.CreateUser(db, user); err != nil {
		// Jika gagal membuat user baru di database, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to create user in database",
		})
	}
	// Mengembalikan response jika berhasil membuat user baru
	return c.JSON(http.StatusCreated, echo.Map{
		"status": "success",
		"data":   user,
	})
}

// TODO: mengupdate user berdasarkan id
func UpdateUser(c echo.Context) error {
	//
	db := c.Get("db").(*gorm.DB)         // Mendapatkan database dari context
	user := new(models.User)             // Membuat variabel baru bertipe models.User
	if err := c.Bind(user); err != nil { // Mengecek apakah ada error saat binding data
		return c.JSON(http.StatusBadRequest, echo.Map{ // Mengembalikan response jika terjadi error
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	if err := models.CreateUser(db, user); err != nil { // Mencoba membuat user baru di database
		return c.JSON(http.StatusInternalServerError, echo.Map{ // Mengembalikan response jika terjadi error
			"status":  "error",
			"message": "Failed to create user in database",
		})
	}
	return c.JSON(http.StatusCreated, echo.Map{ // Mengembalikan response jika berhasil membuat user baru
		"status": "success",
		"data":   user,
	})
}
