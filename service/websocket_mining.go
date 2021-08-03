package service

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
)

const (
	miningSetDifficultyKey = "mining.set_difficulty"
	miningNotifyKey        = "mining.notify"
)

func (ws *webSocket) handleMiningAuthorize(req *rpcRequest) {
	log.Print("[mining.authorize] request")

	var response *rpcResponse
	if ws.isValidMiningAuthorize(req) {
		response = &rpcResponse{ID: req.ID, Result: true}
	} else {
		response = &rpcResponse{Error: errRPCInvalidParams}
	}

	ws.WriteMsg(response)
}

func (ws *webSocket) handleMiningSubscribe(req *rpcRequest) {
	log.Print("[mining.subscribe] request")

	var response *rpcResponse
	if ws.subscription != nil {
		log.Print("already subscribed!")
		response = &rpcResponse{Error: errRPCInvalidParams}
	} else {
		if ws.isRequestingExistingSubscription(req) {
			response = ws.handleExistingSubscription(req)
		} else {
			response = ws.handleNewSubscription(req)
		}
	}

	ws.WriteMsg(response)
}

func (ws *webSocket) handleExistingSubscription(req *rpcRequest) *rpcResponse {
	subscriber := req.Params[0]
	extraNonce1, err := strconv.ParseInt(req.Params[1], 16, 64)
	if err != nil {
		log.Printf("error converting hexadecimal extraNonce1 to integer value: %v", err)
		return &rpcResponse{Error: errRPCInvalidParams}
	}

	subscription, err := ws.svc.getExistingSubscription(subscriber, extraNonce1)
	if err != nil {
		return &rpcResponse{Error: errRPCInternal}
	}
	if subscription == nil {
		return &rpcResponse{Error: errRPCInvalidParams}
	}
	if subscription.activeSession {
		log.Printf("subscription from subscriber: %s and extraNonce1: %d is already active", subscription.subscriber, extraNonce1)
		return &rpcResponse{Error: errRPCInvalidParams}
	}

	ws.subscription = subscription
	return ws.buildSubscriptionRPCResponse(req.ID, subscription)
}

func (ws *webSocket) handleNewSubscription(req *rpcRequest) *rpcResponse {
	var subscriber string
	if len(req.Params) > 0 && req.Params[0] != "" {
		subscriber = req.Params[0]
	} else {
		// if no param is received, random uuid is assigned
		subscriber = uuid.NewString()
	}

	subscription, err := ws.svc.createSubscription(subscriber)
	if err != nil {
		return &rpcResponse{Error: errRPCInternal}
	}

	ws.subscription = subscription
	return ws.buildSubscriptionRPCResponse(req.ID, subscription)
}

func (ws *webSocket) buildSubscriptionRPCResponse(requestId int64, subscription *subscription) *rpcResponse {
	return &rpcResponse{ID: requestId, Result: []interface{}{
		[]interface{}{
			[]string{miningSetDifficultyKey, subscription.setDifficulty},
			[]string{miningNotifyKey, subscription.notify},
		},
		fmt.Sprintf("%08x", subscription.extraNonce1),
		subscription.extraNonce2,
	}}
}

func (ws *webSocket) isValidMiningAuthorize(req *rpcRequest) bool {
	if len(req.Params) != 2 {
		return false
	}
	if req.Params[0] == "" {
		return false
	}
	return true
}

func (ws *webSocket) isRequestingExistingSubscription(req *rpcRequest) bool {
	return len(req.Params) == 2 && req.Params[0] != "" && req.Params[1] != ""
}
