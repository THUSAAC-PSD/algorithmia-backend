package requestemailverification

import (
	"fmt"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*user.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *user.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
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

		if err := ctx.Validate(command); err != nil {
			return err
		}

		code, err := e.handler.Handle(ctx.Request().Context(), command)
		if errors.Is(err, ErrEmailTimedOut) {
			return httperror.New(http.StatusTooManyRequests, fmt.Sprintf("You can only send one email every %d minutes", timeoutDurationMins)).
				WithType(httperror.ErrTypeRateLimitExceeded)
		} else if errors.Is(err, ErrEmailAssociatedWithUser) {
			return httperror.New(http.StatusUnprocessableEntity, "This email is already associated with an existing user").
				WithType(httperror.ErrTypeUserAlreadyExists)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		// In development mode (when email verification is disabled), return the code
		if !e.handler.requireEmailVerification && code != "" {
			return ctx.JSON(http.StatusOK, map[string]string{
				"message": "Email verification disabled (development mode)",
				"code":    code,
			})
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
