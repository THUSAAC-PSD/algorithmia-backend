package manageuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type ResetPasswordCommandHandler struct {
	repo           Repository
	passwordHasher PasswordHasher
	authProvider   contract.AuthProvider
	validator      *validator.Validate
	uowFactory     contract.UnitOfWorkFactory
	l              logger.Logger
}

func NewResetPasswordCommandHandler(
	repo Repository,
	passwordHasher PasswordHasher,
	authProvider contract.AuthProvider,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *ResetPasswordCommandHandler {
	return &ResetPasswordCommandHandler{
		repo:           repo,
		passwordHasher: passwordHasher,
		authProvider:   authProvider,
		validator:      validator,
		uowFactory:     uowFactory,
		l:              l,
	}
}

func (h *ResetPasswordCommandHandler) Handle(ctx context.Context, command *ResetPasswordCommand) error {
	if command == nil {
		return errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.StructCtx(ctx, command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	can, err := h.authProvider.Can(ctx, constant.PermissionUserManageRolesAny)
	if err != nil {
		return errors.WrapIf(err, "failed to check permission for resetting password")
	}
	if !can {
		return customerror.NewNoPermissionError(constant.PermissionUserManageRolesAny)
	}

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		user, err := h.repo.GetUserWithRoles(ctx, command.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WithStack(ErrUserNotFound)
			}
			return errors.WrapIf(err, "failed to load user")
		}

		hashedPassword, err := h.passwordHasher.Hash(command.NewPassword)
		if err != nil {
			return errors.WrapIf(err, "failed to hash password")
		}

		if err := h.repo.UpdatePassword(ctx, user.UserID, hashedPassword); err != nil {
			return errors.WrapIf(err, "failed to update password")
		}

		return nil
	})
}
