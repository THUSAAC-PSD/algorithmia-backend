package deleteproblemdraft

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

func (r *GormRepository) DeleteProblemDraft(ctx context.Context, problemDraftID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	if res := db.WithContext(ctx).
		Delete(&database.ProblemDraft{
			ProblemDraftID: problemDraftID,
		}); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to delete problem draft")
	} else if res.RowsAffected == 0 {
		return ErrProblemDraftNotFound
	}

	return nil
}
