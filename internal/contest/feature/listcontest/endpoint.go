package listcontest

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
	e.ContestsGroup.GET("", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		response, err := mediatr.Send[*Query, *Response](
			ctx.Request().Context(),
			&Query{},
		)
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusOK, response)
	}
}
