package assignproblem

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
	ErrProblemNotFound = errors.New("problem not found")
	ErrContestNotFound = errors.New("contest not found")
	ErrTooManyProblems = errors.New("too many problems")
)

type Repository interface {
	DoesContestExist(ctx context.Context, contestID uuid.UUID) (bool, error)
	DoesProblemExist(ctx context.Context, problemID uuid.UUID) (bool, error)
	AssignProblemToContest(ctx context.Context, problemID uuid.UUID, contestID uuid.UUID) error
	IsContestAlmostMaxedOut(ctx context.Context, contestID uuid.UUID) (bool, error)
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

	// TODO: check if the current user has permissions or not

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		if ok, err := h.repo.DoesContestExist(ctx, command.ContestID); err != nil {
			return errors.WrapIf(err, "failed to check if contest exists")
		} else if !ok {
			return errors.WithStack(ErrContestNotFound)
		}

		if ok, err := h.repo.DoesProblemExist(ctx, command.ProblemID); err != nil {
			return errors.WrapIf(err, "failed to check if problem exists")
		} else if !ok {
			return errors.WithStack(ErrProblemNotFound)
		}

		if notOk, err := h.repo.IsContestAlmostMaxedOut(ctx, command.ContestID); err != nil {
			return errors.WrapIf(err, "failed to get contest problem count")
		} else if notOk {
			return errors.WithStack(ErrTooManyProblems)
		}

		if err := h.repo.AssignProblemToContest(ctx, command.ProblemID, command.ContestID); err != nil {
			return errors.WrapIf(err, "failed to assign problem to contest")
		}

		return nil
	})
}
