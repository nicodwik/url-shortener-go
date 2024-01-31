package helpers

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ValidationResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func ResponseOK(message string, data interface{}) interface{} {
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return response
}

func ResponseNotFound(message string) interface{} {
	response := Response{
		Status:  "error",
		Message: message,
		Data:    []string{},
	}
	return response
}

func ResponseValidationError(validationErrors []string) interface{} {
	response := ValidationResponse{
		Status:  "error",
		Message: "Validation Error",
		Errors:  validationErrors,
	}
	return response
}
