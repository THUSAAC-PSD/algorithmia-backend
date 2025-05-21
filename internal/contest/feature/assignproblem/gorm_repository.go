package assignproblem

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

func (r *GormRepository) IsContestAlmostMaxedOut(ctx context.Context, contestID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	type resultModel struct {
		ProblemCount    uint64 `gorm:"problem_count"`
		MaxProblemCount uint64 `gorm:"max_problem_count"`
	}

	var data resultModel
	if err := db.WithContext(ctx).
		Table("contests c").
		Select("COUNT(p.*) AS problem_count, c.max_problem_count AS max_problem_count").
		Joins("LEFT JOIN problems p ON p.assigned_contest_id = c.contest_id").
		Where("contest_id = ?", contestID).
		Group("c.contest_id").
		Find(&data).Error; err != nil {
		return false, errors.WrapIf(err, "failed to get contest data")
	}

	return data.ProblemCount+1 >= data.MaxProblemCount, nil
}

func (r *GormRepository) AssignProblemToContest(ctx context.Context, problemID uuid.UUID, contestID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	problem := database.Problem{ProblemID: problemID}

	if err := db.WithContext(ctx).
		Model(&problem).
		Update("assigned_contest_id", contestID).Error; err != nil {
		return errors.WrapIf(err, "failed to set assigned contest id")
	}

	return nil
}
