package getproblem

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type Repository interface {
	GetProblem(ctx context.Context, problemID uuid.UUID) (*ResponseProblem, error)
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

func (q *QueryHandler) Handle(ctx context.Context, query *Query) (*Response, error) {
	if query == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	problem, err := q.repo.GetProblem(ctx, query.ProblemID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get problem")
	}

	return &Response{
		Problem: *problem,
	}, nil
}
