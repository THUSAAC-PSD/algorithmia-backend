package testproblem

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

var ErrProblemNotApprovedForTesting = errors.New("problem not approved for testing")

type Repository interface {
	GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error)
	CreateTestResult(
		ctx context.Context,
		command *Command,
		testerID uuid.UUID,
		versionID uuid.UUID,
		createdAt time.Time,
	) (uuid.UUID, error)
	GetProblem(ctx context.Context, problemID uuid.UUID) (dto.ProblemStatusAndVersion, error)
	UpdateProblemStatus(ctx context.Context, problemID uuid.UUID, status constant.ProblemStatus) error
	SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error
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

		if problem.Status != constant.ProblemStatusApprovedForTesting {
			return nil, errors.WithStack(ErrProblemNotApprovedForTesting)
		}

		versionID, err := h.repo.GetLatestProblemVersionID(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get latest problem version ID")
		}

		resultID, err := h.repo.CreateTestResult(ctx, command, user.UserID, versionID, time.Now())
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create test result")
		}

		var problemStatus constant.ProblemStatus
		switch command.Status {
		case StatusPassed:
			problemStatus = constant.ProblemStatusAwaitingFinalCheck
		case StatusFailed:
			problemStatus = constant.ProblemStatusNeedsRevision
			if err := h.repo.SetProblemDraftActive(ctx, problem.DraftID); err != nil {
				return nil, errors.WrapIf(err, "failed to set problem draft active")
			}
		}

		if err := h.repo.UpdateProblemStatus(ctx, command.ProblemID, problemStatus); err != nil {
			return nil, errors.WrapIf(err, "failed to update problem status")
		}

		return &Response{
			TestResultID:     resultID,
			ProblemVersionID: versionID,
		}, nil
	})
}
