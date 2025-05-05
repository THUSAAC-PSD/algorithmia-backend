package deletecontest

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type Repository interface {
	DeleteContest(ctx context.Context, contestID uuid.UUID) error
}

type CommandHandler struct {
	repo      Repository
	validator *validator.Validate
}

func NewCommandHandler(repo Repository, validator *validator.Validate) *CommandHandler {
	return &CommandHandler{
		repo:      repo,
		validator: validator,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return mediatr.Unit{}, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	contestID, err := uuid.Parse(command.ContestID)
	if err != nil {
		return mediatr.Unit{}, errors.WrapIf(customerror.ErrValidationFailed, "failed to parse contest ID")
	}

	if err := h.repo.DeleteContest(ctx, contestID); err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to delete contest in repository")
	}

	return mediatr.Unit{}, nil
}
