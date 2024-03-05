package helpers

import "log"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
		Status:  "not-found",
		Message: message,
		Data:    []string{},
	}
	return response
}

func ResponseValidationError(validationErrors []string) interface{} {
	response := ValidationResponse{
		Status:  "validation-error",
		Message: "Validation Error",
		Errors:  validationErrors,
	}
	return response
}

func ResponseServerError(message string, err error) interface{} {
	log.Println("server-error: " + err.Error())

	response := Response{
		Status:  "server-error",
		Message: message,
		Error:   err.Error(),
	}
	return response
}

func ResponseUnauthorized(err error) interface{} {
	log.Println("unauthorized: " + err.Error())

	response := Response{
		Status:  "unauthorized",
		Message: "Unauthorized, please re-login!",
		Error:   err.Error(),
	}
	return response
}
