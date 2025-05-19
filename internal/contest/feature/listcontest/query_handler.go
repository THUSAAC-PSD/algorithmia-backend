package listcontest

import (
	"context"
)

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

func (q *QueryHandler) Handle(ctx context.Context) (*Response, error) {
	contests, err := q.repo.GetAllContests(ctx)
	if err != nil {
		return nil, err
	}

	return &Response{
		Contests: contests,
	}, nil
}
