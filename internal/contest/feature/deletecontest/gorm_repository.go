package deletecontest

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
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	contest := &database.Contest{ContestID: contestID}
	if err := db.WithContext(ctx).Delete(&contest).Error; err != nil {
		return errors.WrapIf(err, "failed to delete contest")
	}
	return nil
}
