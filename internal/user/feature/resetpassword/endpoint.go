package resetpassword

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
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
	e.UsersGroup.POST("/reset-password", e.handle())
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

		if err := e.handler.Handle(ctx.Request().Context(), command); err != nil {
			if errors.Is(err, customerror.ErrCommandNil) ||
				errors.Is(err, customerror.ErrValidationFailed) {
				return err
			}

			switch {
			case errors.Is(err, ErrUserNotFound):
				return httperror.New(http.StatusNotFound, "User not found").WithInternal(err)
			default:
				return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
			}
		}

		return ctx.JSON(http.StatusOK, Response{Message: "Password updated successfully"})
	}
}
