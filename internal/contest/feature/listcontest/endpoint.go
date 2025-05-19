package listcontest

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*contest.EndpointParams
	handler *QueryHandler
}

func NewEndpoint(params *contest.EndpointParams, handler *QueryHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ContestsGroup.GET("", e.handle())
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
