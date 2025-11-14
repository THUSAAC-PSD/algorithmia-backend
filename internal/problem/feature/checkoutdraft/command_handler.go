package checkoutdraft

import (
	"context"
	"time"

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
	ErrProblemNotFound = errors.New("problem not found")
	ErrForbidden       = errors.New("forbidden to edit this problem")
)

type ProblemSummary struct {
	ProblemID      uuid.UUID
	ProblemDraftID uuid.UUID
	CreatorID      uuid.UUID
}

type Repository interface {
	GetProblemSummary(ctx context.Context, problemID uuid.UUID) (ProblemSummary, error)
	GetLatestVersion(ctx context.Context, problemID uuid.UUID) (*VersionAggregate, error)
	ReplaceDraftFromVersion(ctx context.Context, draftID uuid.UUID, version *VersionAggregate, updatedAt time.Time) error
	GetProblemDraft(ctx context.Context, draftID uuid.UUID) (*dto.ProblemDraft, error)
}

type VersionAggregate struct {
	ProblemVersionID    uuid.UUID
	ProblemDifficultyID uuid.UUID
	Details             []VersionDetail
	Examples            []VersionExample
}

type VersionDetail struct {
	Language     string
	Title        string
	Background   string
	Statement    string
	InputFormat  string
	OutputFormat string
	Note         string
}

type VersionExample struct {
	Input  string
	Output string
}

type CommandHandler struct {
	repo         Repository
	authProvider contract.AuthProvider
	validator    *validator.Validate
	uowFactory   contract.UnitOfWorkFactory
	l            logger.Logger
}

func NewCommandHandler(
	repo Repository,
	authProvider contract.AuthProvider,
	validator *validator.Validate,
	uowFactory contract.UnitOfWorkFactory,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:         repo,
		authProvider: authProvider,
		validator:    validator,
		uowFactory:   uowFactory,
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
		return nil, errors.WrapIf(err, "failed to get user")
	}

	uow := h.uowFactory.New()
	result, err := uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (*dto.ProblemDraft, error) {
		summary, err := h.repo.GetProblemSummary(ctx, command.ProblemID)
		if err != nil {
			if errors.Is(err, ErrProblemNotFound) {
				return nil, err
			}
			return nil, errors.WrapIf(err, "failed to get problem summary")
		}

		if summary.CreatorID != user.UserID {
			return nil, errors.WithStack(ErrForbidden)
		}

		version, err := h.repo.GetLatestVersion(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get latest problem version")
		}

		if err := h.repo.ReplaceDraftFromVersion(ctx, summary.ProblemDraftID, version, time.Now()); err != nil {
			return nil, errors.WrapIf(err, "failed to synchronize draft from version")
		}

		draft, err := h.repo.GetProblemDraft(ctx, summary.ProblemDraftID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to fetch problem draft")
		}

		return draft, nil
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		ProblemDraft: *result,
	}, nil
}
