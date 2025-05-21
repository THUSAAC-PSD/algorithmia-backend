package listtester

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"emperror.dev/errors"
)

type Repository interface {
	GetTesters(ctx context.Context) ([]ResponseTester, error)
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
	_, err := q.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user ID from auth provider")
	}

	testers, err := q.repo.GetTesters(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get testers")
	}

	return &Response{
		Testers: testers,
	}, nil
}
