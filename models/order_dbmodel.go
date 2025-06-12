package models

type OrderDbModel struct {
	OrderID    string      `json:"order_id"`
	TableID    string      `json:"table_id"`
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `json:"items"`
	CreatedAt  string
}
