package requestemailverification

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
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
	CreateEmailVerificationCode(ctx context.Context, email string, code string) error
}

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, email string, code string) error
}

type CommandHandler struct {
	repo        Repository
	emailSender EmailSender
	validator   *validator.Validate
	uowFactory  contract.UnitOfWorkFactory
	l           logger.Logger
}

func NewCommandHandler(
	repo Repository,
	emailSender EmailSender,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:        repo,
		emailSender: emailSender,
		validator:   validator,
		uowFactory:  uowFactory,
		l:           l,
	}
}

func (c *CommandHandler) Handle(
	ctx context.Context,
	command *Command,
) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := c.validator.StructCtx(ctx, command); err != nil {
		return mediatr.Unit{}, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	uow := c.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, c.l, func(ctx context.Context) (mediatr.Unit, error) {
		if ok, err := c.repo.IsNotTimedOut(ctx, command.Email); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to check if email is not timed out")
		} else if !ok {
			return mediatr.Unit{}, errors.WithStack(ErrEmailTimedOut)
		}

		if ok, err := c.repo.IsNotAssociatedWithUser(ctx, command.Email); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to check if email is associated with user")
		} else if !ok {
			return mediatr.Unit{}, errors.WithStack(ErrEmailAssociatedWithUser)
		}

		code, err := c.generateVerificationCode()
		if err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to generate verification code")
		}

		if err := c.repo.CreateEmailVerificationCode(ctx, command.Email, code); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to create email verification code")
		}

		if err := c.emailSender.SendVerificationEmail(ctx, command.Email, code); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to send verification email")
		}

		return mediatr.Unit{}, nil
	})
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
