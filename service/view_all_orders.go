package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type ViewOrders struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateViewOrders(db db.DbOperations, sessionCache *store.MemoryStore) *ViewOrders {
	return &ViewOrders{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *ViewOrders) DecodeAndValidate(reqBytes []byte, userType, key string) (any, error) {

	var orderRequest models.OrderRequest
	paymentStatus := string(reqBytes)
	if userType == "CUSTOMER" {
		sessionId, tableID := t.SessionCache.GetSessionAndTableId(key)
		orderRequest.TableId = tableID
		orderRequest.SessionID = sessionId
	} else if userType == "ADMIN" {

		orderRequest.PaymentStatus = paymentStatus
	} else {
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}
	orderRequest.UserType = userType

	return orderRequest, nil
}

func (t *ViewOrders) ProcessRequest(req any) ([]byte, error) {
	logger.Log.Infoln("Inside Order Items service")
	viewRequest, ok := req.(models.OrderRequest)
	if !ok {
		return nil, fmt.Errorf("expecting models.Product")
	}
	orders, err := t.Db.GetOrders(viewRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	respBytes, err := json.Marshal(orders)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	logger.Log.Infoln("Exiting Order Items service")
	return respBytes, nil

}
