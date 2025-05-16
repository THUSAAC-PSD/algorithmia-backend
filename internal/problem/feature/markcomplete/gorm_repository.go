package markcomplete

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared/infrastructure"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type GormRepository struct {
	problemActionRepo infrastructure.ProblemActionRepository
}

func NewGormRepository(problemActionRepo infrastructure.ProblemActionRepository) *GormRepository {
	return &GormRepository{problemActionRepo: problemActionRepo}
}

func (g *GormRepository) GetProblemStatus(ctx context.Context, problemID uuid.UUID) (constant.ProblemStatus, error) {
	problem, err := g.problemActionRepo.GetProblem(ctx, problemID)
	if err != nil {
		var emptyStatus constant.ProblemStatus
		return emptyStatus, errors.WrapIf(err, "failed to get problem")
	}

	return problem.Status, nil
}

func (g *GormRepository) MarkProblemCompleted(ctx context.Context, problemID uuid.UUID) error {
	return g.problemActionRepo.UpdateProblemStatus(ctx, problemID, constant.ProblemStatusCompleted)
}
