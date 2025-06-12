package db

import "github.com/jithesh-wq/oolio-food-cart/models"

type DbOperations interface {
	GetProducts(id string) ([]models.Product, error)
	InsertProducts(products []models.Product) error
	UpsertOrder(sessionID, tableID string, items []models.OrderItem) (string, error)
	UpdatePaymentStatus(orderID string) error
	GetOrders(req models.OrderRequest) ([]models.Order, error)
}
