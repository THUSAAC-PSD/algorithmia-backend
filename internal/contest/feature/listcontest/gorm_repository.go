package listcontest

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
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

func (g *GormRepository) GetAllContests(ctx context.Context) ([]Contest, error) {
	db := database.GetDBFromContext(ctx, g.db)

	var contests []Contest
	if err := db.WithContext(ctx).Find(&contests).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all contests")
	}

	return contests, nil
}
