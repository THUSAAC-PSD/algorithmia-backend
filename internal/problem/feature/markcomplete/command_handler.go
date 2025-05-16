package markcomplete

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

var ErrProblemNotAwaitingFinalCheck = errors.New("problem is not awaiting final check")

type Repository interface {
	MarkProblemCompleted(ctx context.Context, problemID uuid.UUID) error
	GetProblemStatus(ctx context.Context, problemID uuid.UUID) (constant.ProblemStatus, error)
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

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return mediatr.Unit{}, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	_, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to get user from auth provider")
	}

	// TODO: check if the current user has permission to mark the problem as completed or not

	uow := h.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (mediatr.Unit, error) {
		if problemStatus, err := h.repo.GetProblemStatus(ctx, command.ProblemID); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to get problem status")
		} else if problemStatus != constant.ProblemStatusAwaitingFinalCheck && problemStatus != constant.ProblemStatusCompleted {
			return mediatr.Unit{}, errors.WithStack(ErrProblemNotAwaitingFinalCheck)
		}

		if err := h.repo.MarkProblemCompleted(ctx, command.ProblemID); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to mark problem as completed")
		}

		return mediatr.Unit{}, nil
	})
}
