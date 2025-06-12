package service

import (
	"encoding/json"
	"fmt"

	"github.com/jithesh-wq/oolio-food-cart/db"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

type TableSessionService struct {
	Db           db.DbOperations
	SessionCache *store.MemoryStore
}

func CreateTableSessionService(db db.DbOperations, sessionCache *store.MemoryStore) *TableSessionService {
	return &TableSessionService{
		Db:           db,
		SessionCache: sessionCache,
	}
}

func (t *TableSessionService) DecodeAndValidate(reqBytes []byte, userType, key string) (any, error) {
	tableID := string(reqBytes)
	//validate the table id provided in request with the tableid master/ for now im just hardcoding this logic
	switch tableID {
	case "T001":
		return "T001", nil
	case "T002":
		return "T002", nil
	default:
		return nil, fmt.Errorf("invalid table id")
	}
}

func (t *TableSessionService) ProcessRequest(req any) ([]byte, error) {
	tableId, ok := req.(string)
	if !ok {
		return nil, fmt.Errorf("expecting string")
	}
	var authResp models.AuthResponse
	session := t.SessionCache.CreateSession(tableId)
	if session == nil {
		authResp.Result = "failure"
		authResp.Remarks = "session already exists for this table"
		authResp.Session = nil
	} else {
		authResp.Result = "success"
		authResp.Remarks = "session created successfully"
		authResp.Session = session
	}
	authRespBytes, err := json.Marshal(authResp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authresp: %w", err)
	}
	return authRespBytes, nil
}
