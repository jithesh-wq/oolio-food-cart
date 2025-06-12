package models

type CheckoutRequest struct {
	OrderId string `json:"order_id"`
	TableID string `json:"table_id"`
}
