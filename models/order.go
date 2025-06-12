package models

import "time"

type OrderRequest struct {
	Items         []Item `json:"items"`
	SessionID     string
	TableId       string
	PaymentStatus string
	UserType      string
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderResponse struct {
	GeneralResponse GeneralResponse `json:"generalResponse"`
	OrderId         string          `json:"id"`
	TotalPrice      float64         `json:"totalPrice"`
	Items           []Item          `json:"items"`
	Products        []Product       `json:"products"`
}

type Order struct {
	ID            string      `json:"id"`
	TableID       string      `json:"table_id"`
	TotalPrice    float64     `json:"total_price"`
	SessionID     string      `json:"session_id"`
	CreatedAt     time.Time   `json:"created_at"`
	Items         []OrderItem `json:"items"`
	PaymentStatus string      `json:"payment_status"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
