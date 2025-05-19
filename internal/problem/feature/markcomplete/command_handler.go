package markcomplete

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var ErrProblemNotAwaitingFinalCheck = errors.New("problem is not awaiting final check")

type Repository interface {
	MarkProblemCompleted(ctx context.Context, problemID uuid.UUID, completerID uuid.UUID, timestamp time.Time) error
	GetProblemStatus(ctx context.Context, problemID uuid.UUID) (constant.ProblemStatus, error)
}

type CommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
	uowFactory   contract.UnitOfWorkFactory
	broadcaster  contract.MessageBroadcaster
	l            logger.Logger
}

func NewCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
	uowFactory contract.UnitOfWorkFactory,
	broadcaster contract.MessageBroadcaster,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
		uowFactory:   uowFactory,
		broadcaster:  broadcaster,
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

	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to get user from auth provider")
	}

	// TODO: check if the current user has permission to mark the problem as completed or not

	uow := h.uowFactory.New()
	return uowhelper.Do(ctx, uow, h.l, func(ctx context.Context) error {
		if problemStatus, err := h.repo.GetProblemStatus(ctx, command.ProblemID); err != nil {
			return errors.WrapIf(err, "failed to get problem status")
		} else if problemStatus != constant.ProblemStatusAwaitingFinalCheck && problemStatus != constant.ProblemStatusCompleted {
			return errors.WithStack(ErrProblemNotAwaitingFinalCheck)
		}

		timestamp := time.Now()
		if err := h.repo.MarkProblemCompleted(ctx, command.ProblemID, user.UserID, timestamp); err != nil {
			return errors.WrapIf(err, "failed to mark problem as completed")
		}

		if err := h.broadcaster.BroadcastCompletedMessage(command.ProblemID, contract.MessageUser{
			UserID:   user.UserID,
			Username: user.Username,
		}, timestamp); err != nil {
			return errors.WrapIf(err, "failed to broadcast completed message")
		}

		return nil
	})
}
