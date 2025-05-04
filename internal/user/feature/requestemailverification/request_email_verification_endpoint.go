package requestemailverification

import (
	"fmt"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"

	"emperror.dev/errors"
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
	e.AuthGroup.POST("/email-verification", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &Command{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, 100, "invalid request")
		}

		if err := e.Validator.StructCtx(ctx.Request().Context(), command); err != nil {
			return httperror.New(http.StatusBadRequest, 100, fmt.Sprintf("validation error: %s", err.Error()))
		}

		_, err := mediatr.Send[*Command, mediatr.Unit](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, ErrEmailTimedOut) {
			return httperror.New(http.StatusTooManyRequests, 103, "email timed out")
		} else if errors.Is(err, ErrEmailAssociatedWithUser) {
			return httperror.New(http.StatusUnprocessableEntity, 102, "email already associated with user")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, 200, fmt.Sprintf("error in sending command: %s", err.Error()))
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
