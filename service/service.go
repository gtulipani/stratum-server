package service

import (
	"context"
	"stratum-server/config"
	"stratum-server/repository"

	"github.com/gorilla/websocket"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

// Service describes service to deal with devices.
type Service interface {
	// Health: returns server status
	Health() *HealthResponse
	// RunWebsocketConnection: creates a ws connection
	RunWebsocketConnection(ctx context.Context, conn *websocket.Conn)

	// GenerateExtraNonce2: creates a valid ExtraNonce2. Right now it returns 4
	GetExtraNonce2() int64
}

type service struct {
	repository         repository.Repository
	subscriptionsTable config.PostgreSQLTableConfig
}

// NewService creates new instance for devices service.
func NewService(repository repository.Repository, subscriptionsTable config.PostgreSQLTableConfig) *service {
	return &service{
		repository:         repository,
		subscriptionsTable: subscriptionsTable,
	}
}
