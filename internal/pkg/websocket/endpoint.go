package websocket

import (
	"context"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	WebsocketGroup *echo.Group
	Hub            *Hub
	AuthProvider   contract.AuthProvider
	Logger         logger.Logger
	Options        *Options
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
	hub *Hub,
	authProvider contract.AuthProvider,
	logger logger.Logger,
	opts *Options,
) *EndpointParams {
	websocketGroup := v1Group.Group.Group("/ws")
	return &EndpointParams{
		WebsocketGroup: websocketGroup,
		Hub:            hub,
		AuthProvider:   authProvider,
		Logger:         logger,
		Options:        opts,
	}
}

type Endpoint struct {
	*EndpointParams
}

func NewEndpoint(ep *EndpointParams) *Endpoint {
	return &Endpoint{
		EndpointParams: ep,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.WebsocketGroup.GET("/chat", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user, err := e.AuthProvider.MustGetUser(ctx.Request().Context())
		if err != nil {
			return httperror.New(http.StatusUnauthorized, "Authentication required for WebSocket").WithInternal(err)
		}

		conn, err := websocket.Accept(ctx.Response().Writer, ctx.Request(), &websocket.AcceptOptions{
			InsecureSkipVerify: e.Options.SkipTLSVerification,
			OriginPatterns:     e.Options.OriginPatterns,
		})
		if err != nil {
			e.Logger.Errorw("WebSocket upgrade failed", map[string]interface{}{
				"error":   err,
				"user_id": user.UserID,
			})

			// websocket.Accept writes its own error response, so just return the error.
			return err
		}

		e.Logger.Infow("WebSocket connection established", map[string]interface{}{
			"user_id": user.UserID,
		})

		client := NewClient(e.Logger, e.Hub, conn, user.UserID)
		e.Hub.register <- client

		clientCtx, cancelClientPumps := context.WithCancel(ctx.Request().Context())

		go client.writePump(clientCtx)
		client.readPump(clientCtx)
		defer cancelClientPumps()

		return nil // Echo handles the response for the successful upgrade
	}
}
