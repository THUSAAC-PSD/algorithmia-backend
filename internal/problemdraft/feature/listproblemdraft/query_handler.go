package listproblemdraft

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/dto"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type Repository interface {
	GetActiveProblemDraftsByCreator(ctx context.Context, userID uuid.UUID) ([]dto.ProblemDraft, error)
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

	problemDrafts, err := q.repo.GetActiveProblemDraftsByCreator(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	return &Response{
		ProblemDrafts: problemDrafts,
	}, nil
}
