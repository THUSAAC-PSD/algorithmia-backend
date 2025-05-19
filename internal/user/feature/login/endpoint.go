package login

import (
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
	e.AuthGroup.POST("/login", e.handle())
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

		err := e.handler.Handle(ctx.Request().Context(), command)
		if errors.Is(err, ErrInvalidCredentials) {
			return httperror.New(http.StatusUnprocessableEntity, "Invalid credentials").
				WithType(httperror.ErrTypeInvalidCredentials)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
