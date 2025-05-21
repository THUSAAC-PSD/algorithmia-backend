package submitproblemdraft

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/dto"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var (
	ErrProblemDraftNotFound      = errors.New("problem draft not found")
	ErrContestNotFound           = errors.New("contest not found")
	ErrProblemDraftNotActive     = errors.New("problem draft is not active")
	ErrProblemDoesntNeedRevision = errors.New("problem does not need revision")
	ErrNotCreator                = errors.New("not the creator of the problem draft")
	ErrMissingProblemDifficulty  = errors.New("problem draft missing difficulty")
)

type Repository interface {
	GetProblemDraft(ctx context.Context, problemDraftID uuid.UUID) (*dto.ProblemDraft, error)
	GetProblemStatus(ctx context.Context, problemID uuid.UUID) (*constant.ProblemStatus, error)
	SetProblemDraftInactive(ctx context.Context, problemDraftID uuid.UUID) error
	UpsertProblemFromDraft(
		ctx context.Context,
		draft *dto.ProblemDraft,
		targetContestID uuid.NullUUID,
		status constant.ProblemStatus,
		updatedAt time.Time,
	) (uuid.UUID, error)
	CreateProblemVersionFromDraft(
		ctx context.Context,
		problemID uuid.UUID,
		draft *dto.ProblemDraft,
		createdAt time.Time,
	) (uuid.UUID, error)
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

	problemDraft, err := h.repo.GetProblemDraft(ctx, command.ProblemDraftID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get problem draft for verification")
	} else if problemDraft == nil {
		return nil, errors.WithStack(ErrProblemDraftNotFound)
	} else if problemDraft.CreatorID != user.UserID {
		return nil, errors.WithStack(ErrNotCreator)
	} else if !problemDraft.IsActive {
		return nil, errors.WithStack(ErrProblemDraftNotActive)
	}

	if problemDraft.ProblemDifficulty.ProblemDifficultyID == uuid.Nil {
		return nil, errors.WithStack(ErrMissingProblemDifficulty)
	}

	if problemDraft.SubmittedProblemID.Valid {
		status, err := h.repo.GetProblemStatus(ctx, problemDraft.SubmittedProblemID.UUID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get problem status")
		} else if status != nil && *status != constant.ProblemStatusNeedsRevision {
			// If the status is nil, it means the problem is not yet created, so we don't need to check the status.
			return nil, errors.WithStack(ErrProblemDoesntNeedRevision)
		}
	}

	uow := h.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (*Response, error) {
		if err := h.repo.SetProblemDraftInactive(ctx, command.ProblemDraftID); err != nil {
			return nil, errors.WrapIf(err, "failed to set problem draft inactive")
		}

		timestamp := time.Now()

		problemID, err := h.repo.UpsertProblemFromDraft(
			ctx,
			problemDraft,
			command.TargetContestID,
			constant.ProblemStatusPendingReview,
			timestamp,
		)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create problem from draft")
		}

		problemVersionID, err := h.repo.CreateProblemVersionFromDraft(ctx, problemID, problemDraft, timestamp)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create problem version from draft")
		}

		details, err := h.authProvider.MustGetUserDetails(ctx, user.UserID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get user details")
		}

		if err := h.broadcaster.BroadcastSubmittedMessage(problemID, contract.MessageUser{
			UserID:   user.UserID,
			Username: details.Username,
		}, timestamp); err != nil {
			return nil, errors.WrapIf(err, "failed to broadcast submitted message")
		}

		return &Response{
			ProblemID:        problemID,
			ProblemVersionID: problemVersionID,
		}, nil
	})
}
