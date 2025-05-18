package websocket

import (
	"context"
	"encoding/json"

	"emperror.dev/errors"
)

type (
	SendErrorFunc func(message string, code ErrorCode)
	SendAckFunc   func(message string)
)

type Handler func(ctx context.Context, payload interface{}, sendError SendErrorFunc, sendAck SendAckFunc)

type Router struct {
	actionHandlers map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		actionHandlers: make(map[string]Handler),
	}
}

func (r *Router) RegisterAction(action string, handler Handler) {
	r.actionHandlers[action] = handler
}

func (r *Router) Trigger(
	ctx context.Context,
	action string,
	payload interface{},
	sendError SendErrorFunc,
	sendAck SendAckFunc,
) bool {
	if handler, ok := r.actionHandlers[action]; ok {
		handler(ctx, payload, sendError, sendAck)
		return true
	}

	return false
}

func (r *Router) UnmarshalPayload(payloadData interface{}, target interface{}) error {
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		return errors.WrapIf(err, "failed to marshal payload")
	}

	if err := json.Unmarshal(payloadBytes, target); err != nil {
		return errors.WrapIf(err, "failed to unmarshal payload")
	}

	return nil
}
