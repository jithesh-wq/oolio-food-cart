package utils

import "github.com/google/uuid"

func GenerateAPIKey() string {
	return uuid.New().String()
}
