package controllers

import (
	"cafe/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TODO:  mengambil semua data food
func GetFoods(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengambil semua data makanan dari database
	foods, err := models.GetFoods(db)
	if err != nil {
		// Jika Gagal mengambil data makanan dari database
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mengambil makanan dari database",
		})
	}
	// Berhasil mengambil data makanan dari database
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   foods,
	})
}

// TODO: menambahkan menu di food
func CreateFood(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)         // Mendapatkan database dari context
	food := new(models.Food)             // Membuat variabel baru bertipe models.Food
	if err := c.Bind(food); err != nil { // Mengecek apakah ada error saat binding data
		return c.JSON(http.StatusBadRequest, echo.Map{ // Mengembalikan response jika terjadi error
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	// jika gagal membuat menu makanan di database
	if err := models.CreateFood(db, food); err != nil { // Mencoba membuat menu makanan di database
		return c.JSON(http.StatusInternalServerError, echo.Map{ // Mengembalikan response jika terjadi error
			"status":  "error",
			"message": "Failed to create food in database",
		})
	}
	// Berhasil membuat menu makanan di database
	return c.JSON(http.StatusCreated, echo.Map{ // Mengembalikan response jika berhasil membuat menu makanan
		"status": "success",
		"data":   food,
	})
}

// TODO: mengupdate data food
func UpdateFood(c echo.Context) error {
	// Mendapatkan database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengubah string id menjadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Jika gagal mengubah string ke integer, maka return bad request
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid food ID",
		})
	}
	// Mendapatkan data makanan berdasarkan id
	food, err := models.GetFoodByID(db, uint(id))

	// Jika gagal mendapatkan data makanan, maka return internal server error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to retrieve food from database",
		})
	}

	// Jika tidak ada data makanan yang ditemukan, maka return not found
	if food == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Food not found",
		})
	}
	// Mengikat data makanan dengan request body

	// Jika gagal mengikat data makanan dengan request body, maka return bad request
	if err := c.Bind(food); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	// Menyimpan data makanan di database

	// Jika gagal menyimpan data makanan di database, maka return internal server error
	if err := models.UpdateFood(db, food); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to update food in database",
		})
	}
	// Return status OK dan data makanan

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   food,
	})
}

// TODO:  menghapus food by id
func DeleteFood(c echo.Context) error {
	// mendapatkan database dari context
	db := c.Get("db").(*gorm.DB)
	// mengubah string ke integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// jika gagal mengubah string ke integer
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "ID makanan tidak valid",
		})
	}

	// mendapatkan makanan berdasarkan id
	food, err := models.GetFoodByID(db, uint(id))

	// jika gagal mendapatkan makanan dari database
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mengambil makanan dari database",
		})
	}

	// jika makanan tidak ditemukan
	if food == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Makanan tidak ditemukan",
		})
	}

	// jika gagal menghapus makanan dari database menghapus makanan dari database
	if err := models.DeleteFood(db, id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal menghapus makanan dari database",
		})
	}

	// jika berhasil menghapus makanan
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": fmt.Sprintf("Makanan dengan ID %d telah dihapus", id),
	})
}

// TODO: menampilkan menu food berdasarkan id
func GetFoodById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Mengkonversi string ke integer
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  "error",
			Message: "ID makanan tidak valid",
		})
	}
	db := c.Get("db").(*gorm.DB)                  // Mendapatkan instance GORM db dari context
	food, err := models.GetFoodByID(db, uint(id)) // Mengkonversi id ke uint
	// JIKA GAGAL MENDAPATKAN DATA MAKA KEMBALIKAN RESPONSE ERROR
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  "error",
			Message: "Gagal mendapatkan makanan dari database",
		})
	}
	// JIKA DATA TIDAK DITEMUKAN
	if food == nil {
		return c.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "Makanan tidak ditemukan",
		})
	}
	// JIKA DATA DITEMUKAN
	return c.JSON(http.StatusOK, models.Response{
		Status: "success",
		Data:   food,
	})
}
