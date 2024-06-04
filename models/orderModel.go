package models

import "time"

type Order struct {
	Order_id     int             `json:"order_id"`
	User_id      int             `json:"user_id"`
	Total_amount float64         `json:"total_amount"`
	Status       string          `json:"status"`
	Created_at   time.Time       `json:"created_at"`
	OrderDetails []Order_details `json:"order_details"`
}


type Order_details struct {
	OrderDetailID int     `json:"order_detail_id"`
	Order_id      int     `json:"order_id"`
	Product_id    int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
}

type OrderResponse struct {
	Order_id     int       `json:"order_id"`
	User_id      int       `json:"user_id"`
	Total_amount float64   `json:"total_amount"`
	Status       string    `json:"status"`
	Created_at   time.Time `json:"created_at"`
}
