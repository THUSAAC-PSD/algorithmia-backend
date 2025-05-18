package reviewproblem

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared/dto"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var ErrProblemNotPendingReview = errors.New("problem not pending review")

type Problem struct {
	Status  constant.ProblemStatus
	DraftID uuid.UUID
}

type Repository interface {
	GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error)
	CreateReview(
		ctx context.Context,
		command *Command,
		reviewerID uuid.UUID,
		versionID uuid.UUID,
		createdAt time.Time,
	) (uuid.UUID, error)
	GetProblem(ctx context.Context, problemID uuid.UUID) (dto.ProblemStatusAndVersion, error)
	UpdateProblemStatus(ctx context.Context, problemID uuid.UUID, status constant.ProblemStatus) error
	UpdateProblemReviewer(ctx context.Context, problemID uuid.UUID, reviewerID uuid.UUID) error
	SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error
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

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (*Response, error) {
	if command == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return nil, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user from auth provider")
	}

	// TODO: check if assigned to this problem or not

	uow := h.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (*Response, error) {
		problem, err := h.repo.GetProblem(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get problem")
		}

		if problem.Status != constant.ProblemStatusPendingReview {
			return nil, errors.WithStack(ErrProblemNotPendingReview)
		}

		versionID, err := h.repo.GetLatestProblemVersionID(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get latest problem version ID")
		}

		timestamp := time.Now()

		reviewID, err := h.repo.CreateReview(ctx, command, user.UserID, versionID, timestamp)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create review")
		}

		var problemStatus constant.ProblemStatus
		switch command.Decision {
		case DecisionApprove:
			problemStatus = constant.ProblemStatusApprovedForTesting
		case DecisionReject:
			problemStatus = constant.ProblemStatusRejected
		case DecisionNeedsRevision:
			problemStatus = constant.ProblemStatusNeedsRevision
		}

		if err := h.repo.UpdateProblemStatus(ctx, command.ProblemID, problemStatus); err != nil {
			return nil, errors.WrapIf(err, "failed to update problem status")
		}

		if command.Decision == DecisionNeedsRevision {
			if err := h.repo.SetProblemDraftActive(ctx, problem.DraftID); err != nil {
				return nil, errors.WrapIf(err, "failed to set problem draft active")
			}
		}

		if err := h.repo.UpdateProblemReviewer(ctx, command.ProblemID, user.UserID); err != nil {
			return nil, errors.WrapIf(err, "failed to update problem reviewer")
		}

		if err := h.broadcaster.BroadcastReviewedMessage(command.ProblemID, contract.MessageUser{
			UserID:   user.UserID,
			Username: user.Username,
		}, string(command.Decision), timestamp); err != nil {
			return nil, errors.WrapIf(err, "failed to broadcast reviewed message")
		}

		return &Response{
			ReviewID:         reviewID,
			ProblemVersionID: versionID,
		}, nil
	})
}
