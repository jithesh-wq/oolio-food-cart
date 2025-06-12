package models

type AuthResponse struct {
	Result  string   `json:"result"`
	Remarks string   `json:"remarks"`
	Session *Session `json:"session,omitempty"`
}
