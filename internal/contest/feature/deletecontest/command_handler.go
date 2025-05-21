package deletecontest

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

func (h *CommandHandler) Handle(ctx context.Context, command *Command) error {
	if command == nil {
		return errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	if err := h.repo.DeleteContest(ctx, command.ContestID); err != nil {
		return errors.WrapIf(err, "failed to delete contest in repository")
	}

	return nil
}
