package createcontest

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

func (r *GormRepository) CreateContest(ctx context.Context, contest Contest) error {
	db := database.GetDBFromContext(ctx, r.db)

	contestModel := database.Contest{
		ContestID:        contest.ContestID,
		Title:            contest.Title,
		Description:      contest.Description,
		MinProblemCount:  contest.MinProblemCount,
		MaxProblemCount:  contest.MaxProblemCount,
		DeadlineDatetime: contest.DeadlineDatetime,
		CreatedAt:        contest.CreatedAt,
	}

	if err := db.WithContext(ctx).Create(&contestModel).Error; err != nil {
		return errors.WrapIf(err, "failed to create contest")
	}
	return nil
}
