package login

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
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
}

func NewCommandHandler(
	repo Repository,
	passwordChecker PasswordChecker,
	sessionManager SessionManager,
	validator *validator.Validate,
) *CommandHandler {
	return &CommandHandler{
		repo:            repo,
		passwordChecker: passwordChecker,
		sessionManager:  sessionManager,
		validator:       validator,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, customerror.ErrCommandNil
	}

	if err := h.validator.StructCtx(ctx, command); err != nil {
		return mediatr.Unit{}, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	user, err := h.repo.GetUserByUsername(ctx, command.Username)
	if err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to get user by username")
	} else if user == nil {
		// Prevent attackers from knowing if the user exists
		_, _ = h.passwordChecker.Check(someHashedPassword, command.Password)
		return mediatr.Unit{}, errors.WithStack(ErrInvalidCredentials)
	}

	ok, err := h.passwordChecker.Check(user.HashedPassword, command.Password)
	if err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to check password")
	} else if !ok {
		return mediatr.Unit{}, errors.WithStack(ErrInvalidCredentials)
	}

	if err := h.sessionManager.SetUser(ctx, *user); err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to set user in session")
	}

	return mediatr.Unit{}, nil
}
