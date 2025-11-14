package manageuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type DeleteCommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
}

func NewDeleteCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
) *DeleteCommandHandler {
	return &DeleteCommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
	}
}

func (h *DeleteCommandHandler) Handle(ctx context.Context, command *DeleteCommand) error {
	if command == nil {
		return errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	can, err := h.authProvider.Can(ctx, constant.PermissionUserManageRolesAny)
	if err != nil {
		return errors.WrapIf(err, "failed to check permission for deleting user")
	}

	if !can {
		return customerror.NewNoPermissionError(constant.PermissionUserManageRolesAny)
	}

	currentUser, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to get current user")
	}

	if currentUser.UserID == command.UserID {
		return errors.WithStack(ErrCannotDeleteSelf)
	}

	targetUser, err := h.repo.GetUserWithRoles(ctx, command.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithStack(ErrUserNotFound)
		}

		return errors.WrapIf(err, "failed to load target user")
	}

	if userHasRole(targetUser.Roles, "super_admin") {
		count, err := h.repo.CountSuperAdmins(ctx)
		if err != nil {
			return errors.WrapIf(err, "failed to count super admins")
		}

		if count <= 1 {
			return errors.WithStack(ErrCannotDeleteLastSuperAdmin)
		}
	}

	if err := h.repo.DeleteUser(ctx, command.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithStack(ErrUserNotFound)
		}

		return errors.WrapIf(err, "failed to delete user")
	}

	return nil
}
