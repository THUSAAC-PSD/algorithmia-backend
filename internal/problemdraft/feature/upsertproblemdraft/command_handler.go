package upsertproblemdraft

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var (
	ErrInvalidProblemDraftID      = errors.New("invalid problem draft ID")
	ErrInvalidProblemDifficultyID = errors.New("invalid problem difficulty ID")
	ErrNotCreatorOrInactive       = errors.New("not the creator of the problem draft or inactive draft")
)

type Repository interface {
	VerifyActiveProblemDraftCreator(ctx context.Context, problemDraftID uuid.UUID, creatorID uuid.UUID) (bool, error)

	UpsertProblemDraft(
		ctx context.Context,
		command *Command,
		createdAt *time.Time,
		updatedAt time.Time,
		creatorID uuid.UUID,
		exampleIDs []uuid.UUID,
		detailIDs []uuid.UUID,
	) error
}

type CommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
}

func NewCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
) *CommandHandler {
	return &CommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
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

	var (
		createdAt *time.Time
		updatedAt = time.Now()
	)

	if !command.ProblemDraftID.Valid {
		createdAt = &updatedAt
		id, err := uuid.NewV7()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to generate new UUID for problem draft")
		}

		command.ProblemDraftID = uuid.NullUUID{UUID: id, Valid: true}
	} else if ok, err := h.repo.VerifyActiveProblemDraftCreator(ctx, command.ProblemDraftID.UUID, user.UserID); err != nil {
		return nil, errors.WrapIf(err, "failed to verify active problem draft creator")
	} else if !ok {
		return nil, errors.WithStack(ErrNotCreatorOrInactive)
	}

	exampleIDs := make([]uuid.UUID, len(command.Examples))
	for i := range command.Examples {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to generate new UUID for example")
		}

		exampleIDs[i] = id
	}

	detailIDs := make([]uuid.UUID, len(command.Details))
	for i := range command.Details {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, errors.WrapIf(err, "failed to generate new UUID for detail")
		}

		detailIDs[i] = id
	}

	if err := h.repo.UpsertProblemDraft(ctx, command, createdAt, updatedAt, user.UserID, exampleIDs, detailIDs); err != nil {
		return nil, errors.WrapIf(err, "failed to upsert problem draft in repository")
	}

	return &Response{
		ProblemDraftID: command.ProblemDraftID.UUID,
	}, nil
}
