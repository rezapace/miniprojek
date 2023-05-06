package models

import "database/sql"

type OrderDetail struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	FoodID   int     `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func CreateOrderDetailInDB(tx *sql.Tx, orderDetail *OrderDetail) error {
	query := "INSERT INTO order_details (order_id, food_id, quantity, price) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(query, orderDetail.OrderID, orderDetail.FoodID, orderDetail.Quantity, orderDetail.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	orderDetail.ID = int(id)
	return nil
}
