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

type FiatCurrencyResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []models.Fiat_currencies `json:"data"`
}

type RequestResponse struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    []models.RequestListModel `json:"data"`
}

type FiatTransactionResponse struct {
	Status  string                               `json:"status"`
	Message string                               `json:"message"`
	Data    []models.Fiat_transaction_list_model `json:"data"`
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

func FiatCurrencies(status string, message string, data []models.Fiat_currencies) FiatCurrencyResponse {
	return FiatCurrencyResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func RequestListResponse(status string, message string, data []models.RequestListModel) RequestResponse {
	return RequestResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func FiatTransactionListResponse(status string, message string, data []models.Fiat_transaction_list_model) FiatTransactionResponse {
	return FiatTransactionResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
