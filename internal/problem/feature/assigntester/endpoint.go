package assigntester

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.ProblemEndpointParams
}

func NewEndpoint(params *shared.ProblemEndpointParams) *Endpoint {
	return &Endpoint{
		ProblemEndpointParams: params,
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

		if err := e.Validator.StructCtx(ctx.Request().Context(), command); err != nil {
			return httperror.New(http.StatusBadRequest, err.Error()).WithInternal(err)
		}

		_, err := mediatr.Send[*Command, mediatr.Unit](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, shared.ErrProblemNotFound) {
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
