package assigntester

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*problem.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *problem.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ProblemsGroup.PUT("/:problem_id/tester", e.handle())
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
		if errors.Is(err, problem.ErrProblemNotFound) {
			return httperror.New(http.StatusNotFound, "The problem does not exist")
		} else if errors.Is(err, ErrTargetUserNotFound) {
			return httperror.New(http.StatusUnprocessableEntity, "The tester you're trying to assign does not exist")
		} else if errors.Is(err, ErrForbiddenToAssignTester) {
			return httperror.New(http.StatusForbidden, "You are not allowed to assign a tester to this problem")
		} else if errors.Is(err, ErrProblemAlreadyCompleted) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem you're trying to assign a tester to is already completed")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
