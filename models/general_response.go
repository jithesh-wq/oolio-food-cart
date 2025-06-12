package models

type GeneralResponse struct {
	Result  string `json:"result"`
	Remarks string `json:"remarks"`
	Error   string `json:"error,omitempty"`
}
