package deletecontest

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*contest.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *contest.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ContestsGroup.DELETE("/:contest_id", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &Command{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request format")
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		err := e.handler.Handle(ctx.Request().Context(), command)
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
