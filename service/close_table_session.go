package service

import (
	"encoding/json"
	"fmt"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type RemoveTableSessionService struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateRemoveTableSessionService(db db.DbOperations, sessionCache *store.MemoryStore) *RemoveTableSessionService {
	return &RemoveTableSessionService{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *RemoveTableSessionService) DecodeAndValidate(reqBytes []byte, userType, key string) (any, error) {
	var checkoutRequest models.CheckoutRequest
	err := json.Unmarshal(reqBytes, &checkoutRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkoutRequest: %w", err)
	}
	return checkoutRequest, nil
}

func (t *RemoveTableSessionService) ProcessRequest(req any) ([]byte, error) {
	checkoutRequest, ok := req.(models.CheckoutRequest)
	if !ok {
		return nil, fmt.Errorf("expecting CheckoutRequest")
	}
	var authResp models.AuthResponse
	err := t.SessionCache.DeleteTableSesion(checkoutRequest.TableID)
	if err != nil {
		authResp.Result = "failure"
		authResp.Remarks = err.Error()
	} else {
		authResp.Result = "success"
		authResp.Remarks = "closed table successfully"
	}
	if err := t.Db.UpdatePaymentStatus(checkoutRequest.OrderId); err != nil {
		authResp.Result = "failure"
		authResp.Remarks = "failed to close table"
	}
	authRespBytes, err := json.Marshal(authResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authresp: %w", err)
	}
	return authRespBytes, nil
}
