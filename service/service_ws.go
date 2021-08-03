package service

import (
	"context"
	"github.com/gorilla/websocket"
)

const (
	defaultExtraNonce2 int64 = 4
)

func (s *service) RunWebsocketConnection(_ context.Context, conn *websocket.Conn) {
	webSocket := NewWebSocket(conn, s)

	// routine to read messages
	go webSocket.Read()
	// routine to write messages
	go webSocket.Write()
	// routine to graceful shutdown
	go webSocket.Shutdown()
}

func (s *service) GetExtraNonce2() int64 {
	return defaultExtraNonce2
}
