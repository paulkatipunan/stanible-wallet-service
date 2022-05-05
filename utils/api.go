package utils

type StdResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func Response(status string, message string, data []string) StdResponse {
	return StdResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
