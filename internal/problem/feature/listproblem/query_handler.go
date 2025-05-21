package listproblem

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllRelatedProblems(
		ctx context.Context,
		userID uuid.UUID,
		showAll bool,
		showCreated bool,
		showAllPendingReview bool,
		showAssignedTesting bool,
	) ([]ResponseProblem, error)
}

type QueryHandler struct {
	repo         Repository
	authProvider contract.AuthProvider
}

func NewQueryHandler(repo Repository, authProvider contract.AuthProvider) *QueryHandler {
	return &QueryHandler{
		repo:         repo,
		authProvider: authProvider,
	}
}

func (q *QueryHandler) Handle(ctx context.Context) (*Response, error) {
	user, err := q.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user ID from auth provider")
	}

	details, err := q.authProvider.MustGetUserDetails(ctx, user.UserID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user details from auth provider")
	}

	var (
		showAll              = false
		showCreated          = false
		showAllPendingReview = false
		showAssignedTesting  = false
	)

	for _, p := range details.Permissions {
		if p == constant.PermissionProblemListAll {
			showAll = true
		}

		if p == constant.PermissionProblemListCreatedOwn {
			showCreated = true
		}

		if p == constant.PermissionProblemListAwaitingReviewAll {
			showAllPendingReview = true
		}

		if p == constant.PermissionProblemListAssignedTest {
			showAssignedTesting = true
		}
	}

	problems, err := q.repo.GetAllRelatedProblems(
		ctx,
		user.UserID,
		showAll,
		showCreated,
		showAllPendingReview,
		showAssignedTesting,
	)
	if err != nil {
		return nil, err
	}

	return &Response{
		Problems: problems,
	}, nil
}
