package testproblem

import (
	"context"
	"slices"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/dto"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var (
	ErrProblemNotPendingTesting = errors.New("problem is not pending testing")
	ErrTesterNotAssigned        = errors.New("tester is not assigned to this problem")
)

type ResultSummary struct {
	TesterID uuid.UUID
	Status   string
}

type Repository interface {
	GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error)
	SaveTestResult(
		ctx context.Context,
		command *Command,
		testerID uuid.UUID,
		versionID uuid.UUID,
		createdAt time.Time,
	) (uuid.UUID, error)
	GetProblem(ctx context.Context, problemID uuid.UUID) (dto.ProblemStatusAndVersion, error)
	UpdateProblemStatus(ctx context.Context, problemID uuid.UUID, status constant.ProblemStatus) error
	SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error
	GetProblemTesterIDs(ctx context.Context, problemID uuid.UUID) ([]uuid.UUID, error)
	GetTestResultsForVersion(ctx context.Context, versionID uuid.UUID) ([]ResultSummary, error)
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

	uow := h.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (*Response, error) {
		problem, err := h.repo.GetProblem(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get problem")
		}

		if problem.Status != constant.ProblemStatusPendingTesting {
			return nil, errors.WithStack(ErrProblemNotPendingTesting)
		}

		testerIDs, err := h.repo.GetProblemTesterIDs(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get problem testers")
		}

		if len(testerIDs) == 0 {
			return nil, errors.WithStack(ErrTesterNotAssigned)
		}

		if !slices.ContainsFunc(testerIDs, func(id uuid.UUID) bool { return id == user.UserID }) {
			return nil, errors.WithStack(ErrTesterNotAssigned)
		}

		versionID, err := h.repo.GetLatestProblemVersionID(ctx, command.ProblemID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get latest problem version ID")
		}

		timestamp := time.Now()

		resultID, err := h.repo.SaveTestResult(ctx, command, user.UserID, versionID, timestamp)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to create test result")
		}

		newStatus := problem.Status
		switch command.Status {
		case StatusPassed:
			results, err := h.repo.GetTestResultsForVersion(ctx, versionID)
			if err != nil {
				return nil, errors.WrapIf(err, "failed to get test results for version")
			}

			allPassed := true
			passedMap := make(map[uuid.UUID]struct{}, len(results))
			for _, result := range results {
				if result.Status == string(StatusFailed) {
					allPassed = false
					break
				}
				if result.Status == string(StatusPassed) {
					passedMap[result.TesterID] = struct{}{}
				}
			}

			if allPassed {
				for _, testerID := range testerIDs {
					if _, ok := passedMap[testerID]; !ok {
						allPassed = false
						break
					}
				}
			}

			if allPassed {
				newStatus = constant.ProblemStatusAwaitingFinalCheck
			} else {
				newStatus = constant.ProblemStatusPendingTesting
			}
		case StatusFailed:
			newStatus = constant.ProblemStatusTestingChangesRequested
			if err := h.repo.SetProblemDraftActive(ctx, problem.DraftID); err != nil {
				return nil, errors.WrapIf(err, "failed to set problem draft active")
			}
		}

		if newStatus != problem.Status {
			if err := h.repo.UpdateProblemStatus(ctx, command.ProblemID, newStatus); err != nil {
				return nil, errors.WrapIf(err, "failed to update problem status")
			}
		}

		details, err := h.authProvider.MustGetUserDetails(ctx, user.UserID)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to get user details")
		}

		if err := h.broadcaster.BroadcastTestedMessage(command.ProblemID, contract.MessageUser{
			UserID:   user.UserID,
			Username: details.Username,
		}, string(command.Status), timestamp); err != nil {
			return nil, errors.WrapIf(err, "failed to broadcast reviewed message")
		}

		return &Response{
			TestResultID:     resultID,
			ProblemVersionID: versionID,
		}, nil
	})
}
