package assignproblem

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"emperror.dev/errors"
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
	e.ContestsGroup.POST("/:contest_id/problems", e.handle())
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
		if errors.Is(err, ErrProblemNotFound) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem does not exist")
		} else if errors.Is(err, ErrContestNotFound) {
			return httperror.New(http.StatusNotFound, "The contest does not exist")
		} else if errors.Is(err, ErrTooManyProblems) {
			return httperror.New(http.StatusUnprocessableEntity, "The contest already has enough problems")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
