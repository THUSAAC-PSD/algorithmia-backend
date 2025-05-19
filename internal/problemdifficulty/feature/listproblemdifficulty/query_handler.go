package listproblemdifficulty

import (
	"context"
)

type Repository interface {
	GetAllProblemDifficulties(ctx context.Context) ([]ProblemDifficulty, error)
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
	contests, err := q.repo.GetAllProblemDifficulties(ctx)
	if err != nil {
		return nil, err
	}

	return &Response{
		ProblemDifficulties: contests,
	}, nil
}
