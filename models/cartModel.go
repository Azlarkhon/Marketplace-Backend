package models

type Cart struct {
	CartID int `json:"cart_id"`
	UserID int `json:"user_id"`
}

type CartProduct struct {
	CartProductID int `json:"cart_product_id"`
	CartID        int `json:"cart_id"`
	ProductID     int `json:"product_id"`
	Quantity      int `json:"quantity"`
}

type CartProductRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartRequest struct {
	UserID int `json:"user_id"`
}
