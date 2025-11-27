package resetpassword

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type CommandHandler struct {
	repo           Repository
	passwordHasher PasswordHasher
	validator      *validator.Validate
	authProvider   contract.AuthProvider
	uowFactory     contract.UnitOfWorkFactory
	l              logger.Logger
}

func NewCommandHandler(
	repo Repository,
	passwordHasher PasswordHasher,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:           repo,
		passwordHasher: passwordHasher,
		validator:      validator,
		authProvider:   authProvider,
		uowFactory:     uowFactory,
		l:              l,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) error {
	if command == nil {
		return errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.StructCtx(ctx, command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to get current user")
	}

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		currentUser, err := h.repo.GetUserByID(ctx, user.UserID)
		if err != nil {
			return errors.WrapIf(err, "failed to load current user")
		}

		if currentUser == nil {
			return errors.WithStack(ErrUserNotFound)
		}

		hashedPassword, err := h.passwordHasher.Hash(command.NewPassword)
		if err != nil {
			return errors.WrapIf(err, "failed to hash new password")
		}

		if err := h.repo.UpdatePassword(ctx, user.UserID, hashedPassword); err != nil {
			return errors.WrapIf(err, "failed to update password")
		}

		return nil
	})
}
