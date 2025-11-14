package listassignedproblems

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetAssignedProblems(ctx context.Context, contestID uuid.UUID) ([]Problem, error)
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
	problems, err := q.repo.GetAssignedProblems(ctx, query.ContestID)
	if err != nil {
		return nil, err
	}

	return &Response{
		Problems: problems,
	}, nil
}
