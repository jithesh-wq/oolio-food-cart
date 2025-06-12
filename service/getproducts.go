package service

import (
	"encoding/json"
	"fmt"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type GetProductsService struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateGetProductsService(db db.DbOperations, sessionCache *store.MemoryStore) *GetProductsService {
	return &GetProductsService{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *GetProductsService) DecodeAndValidate(reqBytes []byte, userType, key string) (any, error) {
	productId := string(reqBytes)
	return productId, nil
}

func (t *GetProductsService) ProcessRequest(req any) ([]byte, error) {
	productId, ok := req.(string)
	if !ok {
		return nil, fmt.Errorf("expecting string")
	}
	logger.Log.Infoln("Processing request for product ID:", productId)
	if productId == "" {
		// If no product ID is provided, return all products
		products, err := t.Db.GetProducts("0")
		if err != nil {
			return nil, fmt.Errorf("notfound")
		}
		if len(products) == 0 {
			logger.Log.Infoln("No products found")
			return nil, fmt.Errorf("notfound")

		}
		respBytes, err := json.Marshal(products)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal products: %w", err)
		}
		return respBytes, nil
	} else {

		products, err := t.Db.GetProducts(productId)
		if err != nil {
			return nil, fmt.Errorf("notfound")
		}
		if len(products) == 0 {
			return nil, fmt.Errorf("notfound")
		}
		respBytes, err := json.Marshal(products[0])
		if err != nil {
			logger.Log.Errorln("Failed to marshal product:", err)
			return nil, fmt.Errorf("notfound")
		}
		return respBytes, nil
	}
}
