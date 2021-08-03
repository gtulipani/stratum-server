package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Send pings to peer with this period.
	pingPeriod = 30 * time.Second

	miningAuthorizeMethod = "mining.authorize"
	miningSubscribeMethod = "mining.subscribe"
)

var (
	// The JSON sent is not a valid Request object.
	errRPCInvalidReq = &rpcError{
		Code:    -32600,
		Message: "Invalid Request",
	}
	// The method does not exist / is not available.
	errRPCMethodNotFound = &rpcError{
		Code:    -32601,
		Message: "Method not found",
	}
	// Invalid Params.
	errRPCInvalidParams = &rpcError{
		Code:    -32602,
		Message: "Invalid params",
	}
	// Internal JSON-RPC error.
	errRPCInternal = &rpcError{
		Code:    -32603,
		Message: "Internal error",
	}
	// An error occurred on the server while parsing the JSON text.
	errRRCParse = &rpcError{
		Code:    -32700,
		Message: "Parse error",
	}

	errInboundMsgDecode = fmt.Errorf("failed to encode incoming message")
	errInboundMsgReq    = fmt.Errorf("invalid rpc request")
)

type rpcRequest struct {
	ID     int64    `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params,omitempty"`
}

type rpcResponse struct {
	ID     int64       `json:"id,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Websocket interface {
	Read()
	Write()
	Shutdown()

	WriteMsg(i interface{})
	CloseConn()
}

type miningConfig struct {
	extraNonce2 int64
}

type webSocket struct {
	svc        *service
	conn       *websocket.Conn
	close      chan struct{}
	inboundMsg chan []byte
	miningConfig
	subscription *subscription
}

func NewWebSocket(
	conn *websocket.Conn,
	svc *service,
) Websocket {
	ws := &webSocket{
		svc:        svc,
		conn:       conn,
		inboundMsg: make(chan []byte, 256),
		close:      make(chan struct{}),
		miningConfig: miningConfig{
			extraNonce2: svc.GetExtraNonce2(),
		},
	}

	return ws
}

func (ws *webSocket) Read() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in Read routine: %v", r)
		}
		ws.CloseConn()
	}()
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Print("unexpected close, shutting down ws")
			}
			break
		}
		ws.handleMessage(message)
	}
}

func (ws *webSocket) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in Write routine: %v", r)
		}
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.inboundMsg:
			if !ok {
				log.Print("channel inboundMsg seems to be closed")
				return
			}

			if err := ws.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to write msg in websocket: %v", err)
				return
			}
		case <-ticker.C:
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("failed to send Ping msg: %v", err)
				return
			}
		}
	}
}

func (ws *webSocket) Shutdown() {
	<-ws.close

	ws.conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing"), time.Now().Add(time.Second*5))
	ws.conn.Close()
	close(ws.inboundMsg)

	if ws.hasActiveSubscription() {
		ws.svc.inactiveSubscription(ws.subscription)
	}

	log.Print("websocket conn ended")
}

func (ws *webSocket) WriteMsg(i interface{}) {
	raw, err := json.Marshal(i)
	if err != nil {
		return
	}
	ws.inboundMsg <- raw
}

func (ws *webSocket) CloseConn() {
	ws.close <- struct{}{}
}

func (ws *webSocket) hasActiveSubscription() bool {
	return ws.subscription != nil
}

func (ws *webSocket) handleMessage(msg []byte) {
	req, err := ws.decodeMessage(msg)
	if err != nil {
		ws.sendError(err)
		return
	}

	switch req.Method {
	case miningAuthorizeMethod:
		ws.handleMiningAuthorize(req)
	case miningSubscribeMethod:
		ws.handleMiningSubscribe(req)
	default:
		ws.WriteMsg(&rpcResponse{ID: req.ID, Error: errRPCMethodNotFound})
		return
	}
}

func (ws *webSocket) decodeMessage(msg []byte) (*rpcRequest, error) {
	req := &rpcRequest{}
	err := json.Unmarshal(msg, req)
	if err != nil {
		return nil, errInboundMsgDecode
	}

	if req.ID == 0 || req.Method == "" {
		return nil, errInboundMsgReq
	}

	return req, nil
}

func (ws *webSocket) sendError(err error) {
	res := ws.buildErrorResponse(err)
	ws.WriteMsg(res)
}

func (ws *webSocket) buildErrorResponse(err error) *rpcResponse {
	var res *rpcResponse
	switch err {
	case nil:
		log.Print("no error")
		res = nil
	case errInboundMsgDecode:
		log.Print("error decoding JSON-RPC message")
		res = &rpcResponse{Error: errRRCParse}
	case errInboundMsgReq:
		log.Print("invalid JSON-RPC message")
		res = &rpcResponse{Error: errRPCInvalidReq}
	default:
		log.Print("input pipe unknown error")
		res = &rpcResponse{Error: errRPCInternal}
	}

	return res
}
