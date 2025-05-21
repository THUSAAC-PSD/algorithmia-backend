package sendmessage

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

type Endpoint struct {
	r         *websocket.Router
	validator *validator.Validate
	handler   *CommandHandler
}

func NewEndpoint(r *websocket.Router, validator *validator.Validate, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		r:         r,
		validator: validator,
		handler:   handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.r.RegisterAction("send-message", e.handle())
}

func (e *Endpoint) handle() websocket.Handler {
	return func(ctx context.Context, payload interface{}, sendError websocket.SendErrorFunc, sendAck websocket.SendAckFunc) {
		var command Command
		if err := e.r.UnmarshalPayload(payload, &command); err != nil {
			sendError("Failed to unmarshal payload", websocket.ErrCodeInvalidPayload)
			return
		}

		if err := e.validator.StructCtx(ctx, command); err != nil {
			sendError("Failed to validate command: "+err.Error(), websocket.ErrCodeInvalidPayload)
			return
		}

		err := e.handler.Handle(ctx, &command)
		if errors.Is(err, problem.ErrProblemNotFound) {
			sendError("The problem does not exist", websocket.ErrCodeInvalidPayload)
			return
		} else if errors.Is(err, ErrMediaNotFound) {
			sendError("The media does not exist", websocket.ErrCodeInvalidPayload)
			return
		} else if errors.Is(err, ErrUserNotPartOfRoom) {
			sendError("The user is not part of the room", websocket.ErrCodeInvalidPayload)
			return
		} else if err != nil {
			sendError("Failed to send message", websocket.ErrCodeInternalServerError)
			return
		}

		sendAck("Message sent successfully")
	}
}
