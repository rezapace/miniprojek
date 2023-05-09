package models

import "gorm.io/gorm"

type OrderDetail struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id"`
	FoodID   uint    `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// todo CreateOrderDetailInDB
func CreateOrderDetailInDB(db *gorm.DB, orderDetail *OrderDetail) error {
	result := db.Create(orderDetail)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
