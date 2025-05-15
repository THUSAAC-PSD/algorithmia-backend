package reviewproblem

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
	e.ProblemsGroup.POST("/:problem_id/reviews", e.handle())
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

		response, err := mediatr.Send[*Command, *Response](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, ErrProblemNotFound) {
			return httperror.New(http.StatusNotFound, "The problem you're trying to review does not exist").
				WithInternal(err)
		} else if errors.Is(err, ErrProblemNotPendingReview) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem you're trying to review is not in a pending review state")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusCreated, response)
	}
}
