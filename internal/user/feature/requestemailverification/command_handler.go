package requestemailverification

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

const (
	emailVerificationCodeLength  = 6
	emailVerificationCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	ErrEmailTimedOut           = errors.New("email timed out")
	ErrEmailAssociatedWithUser = errors.New("email is associated with a user")
)

type Repository interface {
	IsNotTimedOut(ctx context.Context, email string) (bool, error)
	IsNotAssociatedWithUser(ctx context.Context, email string) (bool, error)
	CreateEmailVerificationCode(ctx context.Context, email string, username string, passwordHash string, code string) error
}

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, email string, code string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type CommandHandler struct {
	repo                     Repository
	emailSender              EmailSender
	passwordHasher           PasswordHasher
	validator                *validator.Validate
	uowFactory               contract.UnitOfWorkFactory
	l                        logger.Logger
	requireEmailVerification bool
}

func NewCommandHandler(
	repo Repository,
	emailSender EmailSender,
	passwordHasher PasswordHasher,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
	cfg *config.Config,
) *CommandHandler {
	return &CommandHandler{
		repo:                     repo,
		emailSender:              emailSender,
		passwordHasher:           passwordHasher,
		validator:                validator,
		uowFactory:               uowFactory,
		l:                        l,
		requireEmailVerification: cfg.RequireEmailVerification,
	}
}

func (c *CommandHandler) Handle(
	ctx context.Context,
	command *Command,
) (string, error) {
	if command == nil {
		return "", errors.WithStack(customerror.ErrCommandNil)
	}

	if err := c.validator.StructCtx(ctx, command); err != nil {
		return "", errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	var generatedCode string

	uow := c.uowFactory.New()
	err := uowhelper.Do(ctx, uow, c.l, func(ctx context.Context) error {
		if ok, err := c.repo.IsNotTimedOut(ctx, command.Email); err != nil {
			return errors.WrapIf(err, "failed to check if email is not timed out")
		} else if !ok {
			return errors.WithStack(ErrEmailTimedOut)
		}

		if ok, err := c.repo.IsNotAssociatedWithUser(ctx, command.Email); err != nil {
			return errors.WrapIf(err, "failed to check if email is associated with user")
		} else if !ok {
			return errors.WithStack(ErrEmailAssociatedWithUser)
		}

		// Hash the password before storing
		hashedPassword, err := c.passwordHasher.Hash(command.Password)
		if err != nil {
			return errors.WrapIf(err, "failed to hash password")
		}

		code, err := c.generateVerificationCode()
		if err != nil {
			return errors.WrapIf(err, "failed to generate verification code")
		}

		generatedCode = code

		if err := c.repo.CreateEmailVerificationCode(ctx, command.Email, command.Username, hashedPassword, code); err != nil {
			return errors.WrapIf(err, "failed to create email verification code")
		}

		// Skip email sending if email verification is not required (development mode)
		if !c.requireEmailVerification {
			c.l.Info("Email verification disabled - skipping email send (development mode)")
			return nil
		}

		if err := c.emailSender.SendVerificationEmail(ctx, command.Email, code); err != nil {
			return errors.WrapIf(err, "failed to send verification email")
		}

		return nil
	})

	return generatedCode, err
}

func (c *CommandHandler) generateVerificationCode() (string, error) {
	var builder strings.Builder
	builder.Grow(emailVerificationCodeLength)

	charSetLength := big.NewInt(int64(len(emailVerificationCodeCharset)))

	for i := 0; i < emailVerificationCodeLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", errors.WrapIf(err, "failed to generate random index for code")
		}

		randomChar := emailVerificationCodeCharset[randomIndex.Int64()]
		builder.WriteByte(randomChar)
	}

	return builder.String(), nil
}
