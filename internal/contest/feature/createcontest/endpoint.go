package createcontest

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.ContestEndpointParams
}

func NewEndpoint(params *shared.ContestEndpointParams) *Endpoint {
	return &Endpoint{
		ContestEndpointParams: params,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ContestsGroup.POST("", e.handle())
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

		if errors.Is(err, ErrInvalidProblemCount) {
			return httperror.New(http.StatusUnprocessableEntity, "Problem count must be greater than 0")
		} else if errors.Is(err, ErrProblemCountRangeFlipped) {
			return httperror.New(http.StatusUnprocessableEntity, "Problem count range flipped. Max should be greater than min")
		} else if errors.Is(err, ErrInvalidDeadlineDatetime) {
			return httperror.New(http.StatusUnprocessableEntity, "Deadline must be in the future")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusCreated, response)
	}
}
