package unassignproblem

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) DoesContestExist(ctx context.Context, contestID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.Contest{}).
		Where("contest_id = ?", contestID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if contest exists")
	}

	return count > 0, nil
}

func (r *GormRepository) DoesProblemExist(ctx context.Context, problemID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ?", problemID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if problem exists")
	}

	return count > 0, nil
}

func (r *GormRepository) UnassignProblemFromContest(ctx context.Context, problemID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	problem := database.Problem{ProblemID: problemID}

	if err := db.WithContext(ctx).
		Model(&problem).
		Update("assigned_contest_id", nil).Error; err != nil {
		return errors.WrapIf(err, "failed to reset assigned contest id")
	}

	return nil
}
