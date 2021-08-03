package service

import "net/http"

type HealthResponse struct {
	Status int64 `json:"status,omitempty"`
}

func (s *service) Health() *HealthResponse {
	return &HealthResponse{
		Status: http.StatusOK,
	}
}
