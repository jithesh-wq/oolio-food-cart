package service

import (
	"encoding/json"
	"fmt"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type OrderItems struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateOrderItems(db db.DbOperations, sessionCache *store.MemoryStore) *OrderItems {
	return &OrderItems{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *OrderItems) DecodeAndValidate(reqBytes []byte, apiuserType, key string) (any, error) {

	var orderRequest models.OrderRequest
	err := json.Unmarshal(reqBytes, &orderRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal orderRequest: %w", err)
	}
	sessionId, tableID := t.SessionCache.GetSessionAndTableId(key)
	orderRequest.TableId = tableID
	orderRequest.SessionID = sessionId
	return orderRequest, nil
}

func (t *OrderItems) ProcessRequest(req any) ([]byte, error) {
	logger.Log.Infoln("Inside Order Items service")
	orderRequest, ok := req.(models.OrderRequest)
	if !ok {
		return nil, fmt.Errorf("expecting []models.Product")
	}
	var products []models.Product
	var OrderItems []models.OrderItem
	var totalPrice float64
	for _, v := range orderRequest.Items {
		var OrderItem models.OrderItem

		p, err := t.Db.GetProducts(v.ProductID)
		if err == nil && len(p) == 1 {
			totalPrice += p[0].Price * float64(v.Quantity)
			products = append(products, p[0])
			OrderItem.ProductID = p[0].ID
			OrderItem.Price = p[0].Price
			OrderItem.Quantity = v.Quantity
			OrderItems = append(OrderItems, OrderItem)
		}
	}
	//create/update3 order in db
	orderId, err := t.Db.UpsertOrder(orderRequest.SessionID, orderRequest.TableId, OrderItems)
	if err != nil {
		logger.Log.Errorln("Failed to upsert order:", err)
		return nil, fmt.Errorf("failed to process order: %w", err)
	}
	response := models.GeneralResponse{
		Result:  "success",
		Remarks: "Order Order placed successfully",
	}
	orderResponse := models.OrderResponse{
		OrderId:         orderId,
		Items:           orderRequest.Items,
		TotalPrice:      totalPrice,
		Products:        products,
		GeneralResponse: response,
	}

	respBytes, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	logger.Log.Infoln("Exiting Order Items service")
	return respBytes, nil

}
