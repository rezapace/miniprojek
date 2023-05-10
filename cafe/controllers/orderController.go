package controllers

import (
	"cafe/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

// TODO: mengambil semua data order
func GetOrders(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengambil semua data order dari database
	orders, err := models.GetOrders(db)
	if err != nil {
		// Gagal mengambil data order dari database
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mengambil order dari database",
		})
	}
	// Berhasil mengambil data order dari database
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   orders,
	})
}

// TODO:  menampilkan order berdasarkan id
func GetOrderById(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengubah string id menjadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Jika gagal mengubah string ke integer, maka return bad request
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "ID order tidak valid",
		})
	}
	// Mendapatkan data order berdasarkan id
	order, err := models.GetOrderById(db, uint(id))
	if err != nil {
		// Jika gagal mendapatkan data order, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mendapatkan order dari database",
		})
	}
	if order == nil {
		// Jika tidak ada data order yang ditemukan, maka return not found
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Order tidak ditemukan",
		})
	}
	// Return status OK dan data order
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   order,
	})
}

// TODO:  membuat order baru
func CreateOrder(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Membuat variabel baru bertipe models.Order
	order := new(models.Order)
	// Mengecek apakah ada error saat binding data
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}
	// Mencoba membuat order baru di database
	if err := models.CreateOrder(db, order); err != nil {
		// Jika gagal membuat order baru di database, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Failed to create order in database",
		})
	}
	// Mengembalikan response jika berhasil membuat order baru
	return c.JSON(http.StatusCreated, echo.Map{
		"status": "success",
		"data":   order,
	})
}

// TODO: Mengubah order
func UpdateOrder(c echo.Context) error {
	// Mendapatkan koneksi database dari context
	db := c.Get("db").(*gorm.DB)
	// Mengubah string id menjadi integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Jika gagal mengubah string ke integer, maka return bad request
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "ID order tidak valid",
		})
	}

	// Mencari order dengan ID yang diberikan di database
	order, err := models.GetOrderById(db, uint(id))
	if err != nil {
		// Jika gagal mendapatkan data order dari database, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mendapatkan order dari database",
		})
	}
	if order == nil {
		// Jika tidak ada data order yang ditemukan, maka return not found
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  "error",
			"message": "Order tidak ditemukan",
		})
	}

	// Mengecek apakah ada error saat binding data
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Mengupdate order di database
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Mengupdate data order
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		// Menghapus semua detail order yang terkait dengan order saat ini
		if err := tx.Where("order_id = ?", order.ID).Delete(&models.OrderDetail{}).Error; err != nil {
			return err
		}

		// Memasukkan kembali detail order yang baru
		for _, detail := range order.Details {
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		// Jika gagal mengupdate order di database, maka return internal server error
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "error",
			"message": "Gagal mengupdate order di database",
		})
	}

	// Return status OK dan data order yang berhasil diupdate
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   order,
	})
}
