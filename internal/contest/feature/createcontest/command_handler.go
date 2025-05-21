package createcontest

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

type Repository interface {
	CreateContest(ctx context.Context, contest Contest) error
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

	if can, err := h.authProvider.Can(ctx, constant.PermissionContestCreate); err != nil {
		return nil, errors.WrapIf(err, "failed to check permission")
	} else if !can {
		return nil, customerror.NewNoPermissionError(constant.PermissionContestCreate)
	}

	contest, err := NewContest(
		command.Title,
		command.Description,
		command.MinProblemCount,
		command.MaxProblemCount,
		command.DeadlineDatetime,
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create contest")
	}

	if err := h.repo.CreateContest(ctx, contest); err != nil {
		return nil, errors.WrapIf(err, "failed to create contest in repository")
	}

	return &Response{
		ContestID: contest.ContestID,
	}, nil
}
