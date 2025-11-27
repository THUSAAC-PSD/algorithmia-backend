package manageuser

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
	listHandler      *ListQueryHandler
	updateHandler    *UpdateCommandHandler
	deleteHandler    *DeleteCommandHandler
	resetPassHandler *ResetPasswordCommandHandler
}

func NewEndpoint(
	params *user.EndpointParams,
	listHandler *ListQueryHandler,
	updateHandler *UpdateCommandHandler,
	deleteHandler *DeleteCommandHandler,
	resetPassHandler *ResetPasswordCommandHandler,
) *Endpoint {
	return &Endpoint{
		EndpointParams:   params,
		listHandler:      listHandler,
		updateHandler:    updateHandler,
		deleteHandler:    deleteHandler,
		resetPassHandler: resetPassHandler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.UsersGroup.GET("", e.handleList())
	e.UsersGroup.PUT("/:user_id", e.handleUpdate())
	e.UsersGroup.DELETE("/:user_id", e.handleDelete())
	e.UsersGroup.POST("/:user_id/reset-password", e.handleResetPassword())
}

func (e *Endpoint) handleList() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		response, err := e.listHandler.Handle(ctx.Request().Context())
		if err != nil {
			if errors.Is(err, customerror.ErrBaseNoPermission) {
				return err
			}
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusOK, response)
	}
}

func (e *Endpoint) handleUpdate() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &UpdateCommand{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request body")
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		response, err := e.updateHandler.Handle(ctx.Request().Context(), command)
		if err != nil {
			if errors.Is(err, customerror.ErrBaseNoPermission) ||
				errors.Is(err, customerror.ErrCommandNil) ||
				errors.Is(err, customerror.ErrValidationFailed) {
				return err
			}

			switch {
			case errors.Is(err, ErrUserNotFound):
				return httperror.New(http.StatusNotFound, "User not found").WithInternal(err)
			case errors.Is(err, ErrRolesRequired):
				return httperror.New(http.StatusBadRequest, "At least one role is required").WithInternal(err)
			case errors.Is(err, ErrRoleNotFound):
				return httperror.New(http.StatusUnprocessableEntity, err.Error()).WithInternal(err)
			case errors.Is(err, ErrEmailAlreadyExists):
				return httperror.New(http.StatusConflict, "Email already exists").WithInternal(err)
			case errors.Is(err, ErrUsernameAlreadyExists):
				return httperror.New(http.StatusConflict, "Username already exists").WithInternal(err)
			case errors.Is(err, ErrCannotRemoveOwnSuperAdmin):
				return httperror.New(http.StatusBadRequest, "You cannot remove your own super admin role").WithInternal(err)
			case errors.Is(err, ErrCannotDeleteLastSuperAdmin):
				return httperror.New(http.StatusUnprocessableEntity, "System must retain at least one super admin").WithInternal(err)
			default:
				return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
			}
		}

		return ctx.JSON(http.StatusOK, response)
	}
}

func (e *Endpoint) handleDelete() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &DeleteCommand{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request")
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		if err := e.deleteHandler.Handle(ctx.Request().Context(), command); err != nil {
			if errors.Is(err, customerror.ErrBaseNoPermission) ||
				errors.Is(err, customerror.ErrCommandNil) ||
				errors.Is(err, customerror.ErrValidationFailed) {
				return err
			}

			switch {
			case errors.Is(err, ErrUserNotFound):
				return httperror.New(http.StatusNotFound, "User not found").WithInternal(err)
			case errors.Is(err, ErrCannotDeleteSelf):
				return httperror.New(http.StatusBadRequest, "You cannot delete your own account").WithInternal(err)
			case errors.Is(err, ErrCannotDeleteLastSuperAdmin):
				return httperror.New(http.StatusUnprocessableEntity, "System must retain at least one super admin").WithInternal(err)
			default:
				return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
			}
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}

func (e *Endpoint) handleResetPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &ResetPasswordCommand{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request")
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		if err := e.resetPassHandler.Handle(ctx.Request().Context(), command); err != nil {
			if errors.Is(err, customerror.ErrBaseNoPermission) ||
				errors.Is(err, customerror.ErrCommandNil) ||
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

		return ctx.JSON(http.StatusOK, map[string]string{"message": "Password updated"})
	}
}
