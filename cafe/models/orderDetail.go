package models

type OrderDetail struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	FoodID   int     `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
