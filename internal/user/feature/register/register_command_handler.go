package register

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

var (
	ErrUserAlreadyExists            = errors.New("user already exists")
	ErrInvalidEmailVerificationCode = errors.New("invalid email verification code")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	IsUserUnique(ctx context.Context, username string, email string) (bool, error)
	CheckAndDeleteEmailVerificationCode(ctx context.Context, email string, code string) (bool, error)
}

type CommandHandler struct {
	repo       Repository
	hasher     PasswordHasher
	validator  *validator.Validate
	uowFactory contract.UnitOfWorkFactory
	l          logger.Logger
}

func NewCommandHandler(
	repo Repository,
	hasher PasswordHasher,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:       repo,
		hasher:     hasher,
		validator:  validator,
		uowFactory: uowFactory,
		l:          l,
	}
}

func (c *CommandHandler) Handle(
	ctx context.Context,
	command *Command,
) (*Response, error) {
	if command == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := c.validator.StructCtx(ctx, command); err != nil {
		return nil, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	uow := c.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, c.l, func(ctx context.Context) (*Response, error) {
		ok, err := c.repo.IsUserUnique(ctx, command.Username, command.Email)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to check if user is unique")
		}
		if !ok {
			return nil, errors.WithStack(ErrUserAlreadyExists)
		}

		ok, err = c.repo.CheckAndDeleteEmailVerificationCode(ctx, command.Email, command.EmailVerificationCode)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to check and delete email verification code")
		}
		if !ok {
			return nil, errors.WithStack(ErrInvalidEmailVerificationCode)
		}

		hashedPassword, err := c.hasher.Hash(command.Password)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to hash password")
		}

		user, err := NewUser(command.Username, command.Email, hashedPassword)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to initialize user")
		}

		err = c.repo.CreateUser(ctx, user)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create user")
		}

		response := &Response{
			User: ResponseUser{
				UserID:    user.UserID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		}

		return response, nil
	})
}
