package controller

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ws(t *testing.T) {
	tests := []struct {
		name    string
		svc     *ServiceMock
		wantErr bool
	}{
		{
			name: "ok",
			svc: &ServiceMock{
				RunWebsocketConnectionFunc: func(ctx context.Context, conn *websocket.Conn) {},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.svc)
			s := httptest.NewServer(h)
			defer s.Close()

			// Convert http://127.0.0.1 to ws://127.0.0.1
			u := "ws" + strings.TrimPrefix(s.URL, "http") + wsEndpoint

			// Connect to the server
			ws, rr, err := websocket.DefaultDialer.Dial(u, nil)
			if (rr.StatusCode != http.StatusSwitchingProtocols) {
				t.Errorf("Request error. status = %d, err: %v", rr.StatusCode, err)
			}
			if ws != nil {
				ws.Close()
			}
		})
	}
}
