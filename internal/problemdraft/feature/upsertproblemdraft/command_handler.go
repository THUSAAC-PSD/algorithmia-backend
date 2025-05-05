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
)

type Repository interface {
	UpsertProblemDraft(
		ctx context.Context,
		command *Command,
		createdAt *time.Time,
		updatedAt time.Time,
		creatorID uuid.UUID,
		exampleIDs []uuid.UUID,
		detailIDs []uuid.UUID,
	) (*ResponseProblemDraft, error)
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

	res, err := h.repo.UpsertProblemDraft(ctx, command, createdAt, updatedAt, user.UserID, exampleIDs, detailIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to upsert problem draft in repository")
	}

	return &Response{
		ProblemDraft: *res,
	}, nil
}
