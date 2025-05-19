package assigntesters

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
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

func (r *GormRepository) DoUsersExist(ctx context.Context, userIDs []uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Where("user_id IN ?", userIDs).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if users exist")
	}

	return int(count) == len(userIDs), nil
}

func (r *GormRepository) UpdateProblemTesters(ctx context.Context, problemID uuid.UUID, testerIDs []uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	problem := database.Problem{ProblemID: problemID}

	testers := make([]database.User, 0, len(testerIDs))
	for _, id := range testerIDs {
		testers = append(testers, database.User{UserID: id})
	}

	if err := db.WithContext(ctx).
		Model(&problem).
		Association("Testers").
		Replace(testers); err != nil {
		return errors.WrapIf(err, "failed to replace problem testers")
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
