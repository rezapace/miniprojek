package models

import (
	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status"`
	OrderTime  gorm.DeletedAt `json:"order_time"`
	Details    []OrderDetail  `json:"details,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderForm struct {
	UserID   uint    `json:"user_id" validate:"required"`
	FoodID   uint    `json:"food_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,min=1"`
	Status   string  `json:"status" validate:"required,oneof=proses selesai"`
}

// todo GetOrderById
func CreateOrderInDB(db *gorm.DB, order *Order) error {
	result := db.Create(order)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
