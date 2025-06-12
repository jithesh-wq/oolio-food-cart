package models

type Session struct {
	SessionID string `json:"session_id"`
	TableId   string `json:"table_id"`
	APIKey    string `json:"api_key"`
	Expiry    int64  `json:"expiry"`
	IsActive  bool   `json:"is_active"`
}

func (s *Session) IsValid() bool {
	return s.IsActive && s.Expiry > 0 && s.APIKey != ""
}
