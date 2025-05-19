package listproblemdraft

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/dto"

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

func (g *GormRepository) GetActiveProblemDraftsByCreator(
	ctx context.Context,
	userID uuid.UUID,
) ([]dto.ProblemDraft, error) {
	db := database.GetDBFromContext(ctx, g.db)

	var problemDrafts []database.ProblemDraft
	if err := db.WithContext(ctx).
		Model(&database.ProblemDraft{}).
		Preload("ProblemDifficulty").
		Preload("ProblemDifficulty.DisplayNames").
		Preload("Details").
		Preload("Examples").
		Preload("SubmittedProblem").
		Where("creator_id = ? AND is_active = ?", userID, true).
		Find(&problemDrafts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all problem drafts")
	}

	result := make([]dto.ProblemDraft, 0, len(problemDrafts))
	for _, problemDraftModel := range problemDrafts {
		pd := dto.FromGormProblemDraft(
			problemDraftModel,
			dto.FromGormProblemDifficulty(problemDraftModel.ProblemDifficulty),
		)
		result = append(result, pd)
	}

	return result, nil
}
