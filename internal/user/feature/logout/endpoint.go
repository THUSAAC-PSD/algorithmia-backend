package logout

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*user.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *user.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.AuthGroup.POST("/logout", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := e.handler.Handle(ctx.Request().Context())
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
