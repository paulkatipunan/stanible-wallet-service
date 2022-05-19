package utils

import (
	"api.stanible.com/wallet/models"
)

type StdResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type TxTypeResponse struct {
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
	Data    []models.Transaction_types `json:"data"`
}

func Response(status string, message string, data []string) StdResponse {
	return StdResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func TxType(status string, message string, data []models.Transaction_types) TxTypeResponse {
	return TxTypeResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
