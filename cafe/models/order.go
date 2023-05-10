package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	UserID     uint          `json:"user_id"`
	TotalPrice float64       `json:"total_price"`
	Status     string        `json:"status"`
	OrderTime  time.Time     `json:"order_time"`
	Details    []OrderDetail `json:"details,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderForm struct {
	UserID   uint    `json:"user_id"`
	FoodID   uint    `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
}

// todo GetOrderById
func CreateOrderInDB(db *gorm.DB, order *Order) error {
	result := db.Create(order)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// mengambil semua data order dari database
func GetOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// menampilkan order berdasarkan id
func GetOrderById(db *gorm.DB, id uint) (*Order, error) {
	var order Order
	if err := db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// membuat order baru di database
func CreateOrder(db *gorm.DB, order *Order) error {
	if err := db.Create(&order).Error; err != nil {
		return err
	}
	return nil
}

// mengupdate order di database
func UpdateOrder(db *gorm.DB, order *Order) error {
	if err := db.Save(&order).Error; err != nil {
		return err
	}
	return nil
}

// menghapus order dari database
func DeleteOrder(db *gorm.DB, order *Order) error {
	if err := db.Delete(&order).Error; err != nil {
		return err
	}
	return nil
}
