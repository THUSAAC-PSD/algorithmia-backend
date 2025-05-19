package markcomplete

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/infrastructure"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	problemActionRepo infrastructure.ProblemActionRepository
	db                *gorm.DB
}

func NewGormRepository(problemActionRepo infrastructure.ProblemActionRepository, db *gorm.DB) *GormRepository {
	return &GormRepository{problemActionRepo: problemActionRepo, db: db}
}

func (g *GormRepository) GetProblemStatus(ctx context.Context, problemID uuid.UUID) (constant.ProblemStatus, error) {
	p, err := g.problemActionRepo.GetProblem(ctx, problemID)
	if err != nil {
		var emptyStatus constant.ProblemStatus
		return emptyStatus, errors.WrapIf(err, "failed to get problem")
	}

	return p.Status, nil
}

func (g *GormRepository) MarkProblemCompleted(
	ctx context.Context,
	problemID uuid.UUID,
	completerID uuid.UUID,
	timestamp time.Time,
) error {
	db := database.GetDBFromContext(ctx, g.db)

	if res := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ?", problemID).
		Update("status", constant.ProblemStatusCompleted).
		Update("completed_at", timestamp).
		Update("completed_by", completerID); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to update problem status")
	} else if res.RowsAffected == 0 {
		return errors.WithStack(problem.ErrProblemNotFound)
	}

	return nil
}
