package service

import (
	"encoding/json"
	"fmt"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type StoreProducts struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateStoreProducts(db db.DbOperations, sessionCache *store.MemoryStore) *StoreProducts {
	return &StoreProducts{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *StoreProducts) DecodeAndValidate(reqBytes []byte, userType, key string) (any, error) {

	var products []models.Product
	err := json.Unmarshal(reqBytes, &products)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal products: %w", err)
	}
	if len(products) == 0 {
		return nil, fmt.Errorf("no products provided")
	}
	for _, product := range products {
		if product.Name == "" || product.Price <= 0 || product.Category == "" ||
			product.Image.Thumbnail == "" || product.Image.Mobile == "" ||
			product.Image.Desktop == "" {
			return nil, fmt.Errorf("invalid product data: %v", product)
		}
	}
	return products, nil
}

func (t *StoreProducts) ProcessRequest(req any) ([]byte, error) {
	logger.Log.Infoln("Inside store products service")
	products, ok := req.([]models.Product)
	if !ok {
		return nil, fmt.Errorf("expecting []models.Product")
	}
	var response models.GeneralResponse
	err := t.Db.InsertProducts(products)
	if err != nil {
		response = models.GeneralResponse{
			Result:  "failed",
			Remarks: "failed to store products",
			Error:   err.Error(),
		}
	} else {
		response = models.GeneralResponse{
			Result:  "success",
			Remarks: "products stored successfully",
		}
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	logger.Log.Infoln("Exiting store products service")
	return respBytes, nil

}
