package assigntester

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var (
	ErrProblemAlreadyCompleted = errors.New("problem already completed")
	ErrTargetUserNotFound      = errors.New("target user not found")
	ErrForbiddenToAssignTester = errors.New("forbidden to assign tester")
)

type Repository interface {
	UpdateProblemTester(ctx context.Context, problemID uuid.UUID, testerID uuid.UUID) error
	IsProblemCompleted(ctx context.Context, problemID uuid.UUID) (bool, error)
	DoesUserExist(ctx context.Context, userID uuid.UUID) (bool, error)
}

type CommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
	uowFactory   contract.UnitOfWorkFactory
	l            logger.Logger
}

func NewCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
		uowFactory:   uowFactory,
		l:            l,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) error {
	if command == nil {
		return errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	_, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to get user from auth provider")
	}

	// TODO: check if the current user has permission to assign tester or not

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		if ok, err := h.repo.DoesUserExist(ctx, command.UserID); err != nil {
			return errors.WrapIf(err, "failed to check if user exists")
		} else if !ok {
			return errors.WithStack(ErrTargetUserNotFound)
		}

		if isCompleted, err := h.repo.IsProblemCompleted(ctx, command.ProblemID); err != nil {
			return errors.WrapIf(err, "failed to check if problem is completed")
		} else if isCompleted {
			return errors.WithStack(ErrProblemAlreadyCompleted)
		}

		if err := h.repo.UpdateProblemTester(ctx, command.ProblemID, command.UserID); err != nil {
			return errors.WrapIf(err, "failed to update problem tester")
		}

		return nil
	})
}
