package models

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

type OrderDetail struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	FoodID   int     `json:"food_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
