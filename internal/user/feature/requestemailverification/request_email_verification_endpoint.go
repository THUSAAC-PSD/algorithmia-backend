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
			return httperror.New(http.StatusBadRequest, "Invalid request format")
		}

		if err := e.Validator.StructCtx(ctx.Request().Context(), command); err != nil {
			return httperror.New(http.StatusBadRequest, err.Error()).WithInternal(err)
		}

		_, err := mediatr.Send[*Command, mediatr.Unit](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, ErrEmailTimedOut) {
			return httperror.New(http.StatusTooManyRequests, fmt.Sprintf("You can only send one email every %d minutes", timeoutDurationMins)).
				WithType(httperror.ErrTypeRateLimitExceeded)
		} else if errors.Is(err, ErrEmailAssociatedWithUser) {
			return httperror.New(http.StatusUnprocessableEntity, "This email is already associated with an existing user").
				WithType(httperror.ErrTypeUserAlreadyExists)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
