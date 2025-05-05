package listproblemdraft

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.ProblemDraftEndpointParams
}

func NewEndpoint(params *shared.ProblemDraftEndpointParams) *Endpoint {
	return &Endpoint{
		ProblemDraftEndpointParams: params,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ProblemDraftsGroup.GET("", e.handle())
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
