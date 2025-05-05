package logout

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.UserEndpointParams
}

func NewEndpoint(params *shared.UserEndpointParams) *Endpoint {
	return &Endpoint{
		UserEndpointParams: params,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.AuthGroup.POST("/logout", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, err := mediatr.Send[*Command, mediatr.Unit](
			ctx.Request().Context(),
			&Command{},
		)
		if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}
