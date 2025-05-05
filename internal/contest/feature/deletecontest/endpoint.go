package deletecontest

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

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
	e.ContestsGroup.DELETE("/:contest_id", e.handle())
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
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
