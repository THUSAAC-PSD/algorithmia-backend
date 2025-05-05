package listcontest

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
)

type Query struct{}

type Repository interface {
	GetAllContests(ctx context.Context) ([]Contest, error)
}

type QueryHandler struct {
	repo Repository
}

func NewQueryHandler(repo Repository) *QueryHandler {
	return &QueryHandler{
		repo: repo,
	}
}

func (q *QueryHandler) Handle(ctx context.Context, query *Query) (*Response, error) {
	if query == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	contests, err := q.repo.GetAllContests(ctx)
	if err != nil {
		return nil, err
	}

	return &Response{
		Contests: contests,
	}, nil
}
