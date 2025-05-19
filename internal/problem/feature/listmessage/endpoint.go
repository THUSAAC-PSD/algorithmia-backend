package listmessage

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*problem.EndpointParams
	handler *QueryHandler
}

func NewEndpoint(params *problem.EndpointParams, handler *QueryHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ProblemsGroup.GET("/:problem_id/messages", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		query := &Query{}
		if err := ctx.Bind(query); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request format")
		}

		if err := ctx.Validate(query); err != nil {
			return err
		}

		response, err := e.handler.Handle(ctx.Request().Context(), query)
		if errors.Is(err, ErrUserNotPartOfRoom) {
			return httperror.New(http.StatusForbidden, "You are not part of this room")
		} else if errors.Is(err, ErrProblemNotFound) {
			return httperror.New(http.StatusNotFound, "The problem was not found")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusOK, response)
	}
}
