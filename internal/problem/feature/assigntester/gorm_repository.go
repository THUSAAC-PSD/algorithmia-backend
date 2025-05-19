package assigntester

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"

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

func (r *GormRepository) DoesUserExist(ctx context.Context, userID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if user exists")
	}

	return count > 0, nil
}

func (r *GormRepository) UpdateProblemTester(ctx context.Context, problemID uuid.UUID, testerID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	if res := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ?", problemID).
		Update("tester_id", testerID); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to update problem tester")
	} else if res.RowsAffected == 0 {
		return errors.WithStack(problem.ErrProblemNotFound)
	}

	return nil
}

func (r *GormRepository) IsProblemCompleted(ctx context.Context, problemID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ? AND status = ?", problemID, constant.ProblemStatusCompleted).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if problem is completed")
	}

	return count > 0, nil
}
