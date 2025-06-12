package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/models"
	"github.com/jithesh-wq/oolio-food-cart/utils"
)

type MemoryStore struct {
	sessions map[string]*models.Session
	tokens   map[string]string
	mu       sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions: make(map[string]*models.Session),
		tokens:   make(map[string]string),
	}
}

// when the user calls an api to login, invoke this function and generate a token for the user/table
func (s *MemoryStore) CreateSession(tableID string) *models.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.tokens[tableID]
	if exists {
		//check the session is active or not
		session, exists := s.sessions[s.tokens[tableID]]
		if exists && session.IsActive && time.Now().Unix() < session.Expiry {
			return session
		}
		return nil
	}
	sessionId := utils.GenerateAPIKey()
	apiKey := utils.GenerateAPIKey()
	newSession := models.Session{
		SessionID: sessionId,
		TableId:   tableID,
		APIKey:    apiKey,
		Expiry:    time.Now().Add(4 * time.Hour).Unix(),
		IsActive:  true,
	}
	s.sessions[apiKey] = &newSession
	s.tokens[tableID] = apiKey
	return &newSession
}

// check the table has an active session api key passed by the user exis
func (s *MemoryStore) ValidateSession(apiKey string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[apiKey]
	if !exists {
		return false
	}
	if session.IsActive && time.Now().Unix() < session.Expiry {
		return true
	}
	return false
}

func (s *MemoryStore) GetSessionAndTableId(apiKey string) (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[apiKey]
	if !exists {
		return "", ""
	}
	return session.SessionID, session.TableId
}

// when the user logs out(after soing payment) he calls this function or this can be called by the restaurant admin
func (s *MemoryStore) DeleteTableSesion(tableID string) error {
	logger.Log.Infoln("Deleting session for table:", tableID)
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.tokens[tableID]
	if !exists {
		return fmt.Errorf("session does not exist for table %s", tableID)
	}
	delete(s.sessions, s.tokens[tableID])
	delete(s.tokens, tableID)
	return nil
}
