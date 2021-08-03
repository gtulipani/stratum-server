package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"stratum-server/service"
)

type APIError struct {
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
}

func encodeHTTPResponse(w http.ResponseWriter, response interface{}) *service.AppError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return &service.AppError{
			Error:   err,
			Message: "error encoding response",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}

func encodeHTTPError(err *service.AppError, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.Code)

	errorBody := APIError{
		Description: err.Error.Error(),
		Message:     err.Message,
	}
	if err := json.NewEncoder(w).Encode(errorBody); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
