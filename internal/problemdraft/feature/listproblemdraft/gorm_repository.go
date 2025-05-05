package listproblemdraft

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared/dto"

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
	var problemDrafts []database.ProblemDraft
	if err := g.db.WithContext(ctx).
		Model(&database.ProblemDraft{}).
		Preload("ProblemDifficulty").
		Preload("ProblemDifficulty.DisplayNames").
		Preload("Details").
		Preload("Examples").
		Where("creator_id = ? AND is_active = ?", userID, true).
		Find(&problemDrafts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all problem drafts")
	}

	result := make([]dto.ProblemDraft, 0, len(problemDrafts))
	for _, problemDraftModel := range problemDrafts {
		pd := dto.ProblemDraft{
			ProblemDraftID: problemDraftModel.ProblemDraftID,
			ProblemDifficulty: dto.ProblemDifficulty{
				ProblemDifficultyID: problemDraftModel.ProblemDifficulty.ProblemDifficultyID,
				DisplayNames: make(
					[]dto.ProblemDifficultyDisplayName,
					len(problemDraftModel.ProblemDifficulty.DisplayNames),
				),
			},
			CreatorID:          problemDraftModel.CreatorID,
			Details:            make([]dto.ProblemDraftDetail, len(problemDraftModel.Details)),
			Examples:           make([]dto.ProblemDraftExample, len(problemDraftModel.Examples)),
			SubmittedProblemID: uuid.NullUUID{Valid: false}, // TODO: link to Problems
			CreatedAt:          problemDraftModel.CreatedAt,
			UpdatedAt:          problemDraftModel.UpdatedAt,
		}

		for i, detail := range problemDraftModel.Details {
			pd.Details[i] = dto.ProblemDraftDetail{
				Language:     detail.Language,
				Title:        detail.Title,
				Background:   detail.Background,
				Statement:    detail.Statement,
				InputFormat:  detail.InputFormat,
				OutputFormat: detail.OutputFormat,
				Note:         detail.Note,
			}
		}

		for i, example := range problemDraftModel.Examples {
			pd.Examples[i] = dto.ProblemDraftExample{
				Input:  example.Input,
				Output: example.Output,
			}
		}

		for i, displayName := range problemDraftModel.ProblemDifficulty.DisplayNames {
			pd.ProblemDifficulty.DisplayNames[i] = dto.ProblemDifficultyDisplayName{
				Language: displayName.Language,
				Name:     displayName.DisplayName,
			}
		}

		result = append(result, pd)
	}

	return result, nil
}
