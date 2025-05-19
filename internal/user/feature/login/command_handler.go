package login

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

const (
	someHashedPassword = "$argon2id$v=19$m=19,t=2,p=1$b1JNdmtaeDNSc1BfWlpP$vHS6id2sbgaYzqOvScoGfImMkyb7U2XhuqdWYaajWFk"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type PasswordChecker interface {
	Check(hashedPassword, plainPassword string) (bool, error)
}

type SessionManager interface {
	SetUser(ctx context.Context, user User) error
}

type CommandHandler struct {
	validator       *validator.Validate
	repo            Repository
	passwordChecker PasswordChecker
	sessionManager  SessionManager
	uowFactory      contract.UnitOfWorkFactory
	l               logger.Logger
}

func NewCommandHandler(
	repo Repository,
	passwordChecker PasswordChecker,
	sessionManager SessionManager,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:            repo,
		passwordChecker: passwordChecker,
		sessionManager:  sessionManager,
		validator:       validator,
		uowFactory:      uowFactory,
		l:               l,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) error {
	if command == nil {
		return customerror.ErrCommandNil
	}

	if err := h.validator.StructCtx(ctx, command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		user, err := h.repo.GetUserByUsername(ctx, command.Username)
		if err != nil {
			return errors.WrapIf(err, "failed to get user by username")
		} else if user == nil {
			// Prevent attackers from knowing if the user exists
			_, _ = h.passwordChecker.Check(someHashedPassword, command.Password)
			return errors.WithStack(ErrInvalidCredentials)
		}

		ok, err := h.passwordChecker.Check(user.HashedPassword, command.Password)
		if err != nil {
			return errors.WrapIf(err, "failed to check password")
		} else if !ok {
			return errors.WithStack(ErrInvalidCredentials)
		}

		if err := h.sessionManager.SetUser(ctx, *user); err != nil {
			return errors.WrapIf(err, "failed to set user in session")
		}

		return nil
	})
}
