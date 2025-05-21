package listtester

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	v1Group *echoweb.V1Group
	handler *QueryHandler
}

func NewEndpoint(v1Group *echoweb.V1Group, handler *QueryHandler) *Endpoint {
	return &Endpoint{
		v1Group: v1Group,
		handler: handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.v1Group.Group.GET("/testers", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		response, err := e.handler.Handle(ctx.Request().Context())
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusOK, response)
	}
}
