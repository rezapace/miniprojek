package models

import "database/sql"

type Order struct {
	ID         int           `json:"id"`
	UserID     int           `json:"user_id"`
	TotalPrice float64       `json:"total_price"`
	Status     string        `json:"status"`
	OrderTime  string        `json:"order_time"`
	Details    []OrderDetail `json:"details,omitempty"`
}

type OrderForm struct {
	UserID   int     `json:"user_id" validate:"required"`
	FoodID   int     `json:"food_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,min=1"`
	Status   string  `json:"status" validate:"required,oneof=proses selesai"`
}

func CreateOrderInDB(tx *sql.Tx, order *Order) error {
	query := "INSERT INTO orders (user_id, total_price, status) VALUES (?, ?, ?)"
	result, err := tx.Exec(query, order.UserID, order.TotalPrice, order.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	order.ID = int(id)
	return nil
}
