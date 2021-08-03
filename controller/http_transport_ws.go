package controller

import (
	"github.com/gorilla/websocket"
	"net/http"
	"stratum-server/service"
)

func ws(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			encodeHTTPError(&service.AppError{
				Error:   err,
				Message: "failed to upgrade connection",
				Code:    http.StatusInternalServerError,
			}, w)
			return
		}
		svc.RunWebsocketConnection(r.Context(), conn)
	}
}
