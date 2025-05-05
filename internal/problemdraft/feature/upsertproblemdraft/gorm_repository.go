package upsertproblemdraft

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) UpsertProblemDraft(
	ctx context.Context,
	command *Command,
	createdAt *time.Time,
	updatedAt time.Time,
	creatorID uuid.UUID,
	exampleIDs []uuid.UUID,
	detailIDs []uuid.UUID,
) (*ResponseProblemDraft, error) {
	db := database.GetDBFromContext(ctx, r.db)

	if !command.ProblemDraftID.Valid {
		return nil, errors.WithStack(ErrInvalidProblemDraftID)
	}

	problemDraftModel := database.ProblemDraft{
		ProblemDraftID:      command.ProblemDraftID.UUID,
		ProblemDifficultyID: command.ProblemDifficultyID,
		CreatorID:           creatorID,
		IsActive:            true,
		Examples:            make([]database.ProblemDraftExample, len(command.Examples)),
		Details:             make([]database.ProblemDraftDetail, len(command.Details)),
		UpdatedAt:           updatedAt,
	}

	if createdAt != nil {
		problemDraftModel.CreatedAt = *createdAt
	}

	for i, example := range command.Examples {
		problemDraftModel.Examples[i] = database.ProblemDraftExample{
			ExampleID:      exampleIDs[i],
			ProblemDraftID: problemDraftModel.ProblemDraftID,
			Input:          example.Input,
			Output:         example.Output,
		}
	}

	for i, detail := range command.Details {
		problemDraftModel.Details[i] = database.ProblemDraftDetail{
			DetailID:       detailIDs[i],
			ProblemDraftID: problemDraftModel.ProblemDraftID,
			Language:       detail.Language,
			Title:          detail.Title,
			Background:     detail.Background,
			Statement:      detail.Statement,
			InputFormat:    detail.InputFormat,
			OutputFormat:   detail.OutputFormat,
			Note:           detail.Note,
		}
	}

	var problemDifficulty database.ProblemDifficulty
	if err := db.WithContext(ctx).
		Model(&database.ProblemDifficulty{}).
		Preload("DisplayNames").
		Where("problem_difficulty_id = ?", command.ProblemDifficultyID).
		First(&problemDifficulty).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrInvalidProblemDifficultyID)
		}

		return nil, errors.WrapIf(err, "failed to get problem difficulty")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("problem_draft_id = ?", command.ProblemDraftID.UUID).Delete(&database.ProblemDraftExample{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete old problem draft examples")
		}

		if err := tx.WithContext(ctx).Where("problem_draft_id = ?", command.ProblemDraftID.UUID).Delete(&database.ProblemDraftDetail{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete old problem draft details")
		}

		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&problemDraftModel).Error; err != nil {
			return errors.WrapIf(err, "failed to upsert problem draft")
		}

		return nil
	}); err != nil {
		return nil, errors.WrapIf(err, "failed to upsert problem draft in transaction")
	}

	response := ResponseProblemDraft{
		ProblemDraftID: problemDraftModel.ProblemDraftID,
		ProblemDifficulty: ResponseProblemDifficulty{
			ProblemDifficultyID: problemDifficulty.ProblemDifficultyID,
			DisplayNames:        make([]ResponseProblemDifficultyDisplayName, len(problemDifficulty.DisplayNames)),
		},
		CreatorID:          creatorID,
		Details:            make([]ResponseProblemDraftDetail, len(problemDraftModel.Details)),
		Examples:           make([]ResponseProblemDraftExample, len(problemDraftModel.Examples)),
		SubmittedProblemID: uuid.NullUUID{Valid: false}, // TODO: link to Problems
		CreatedAt:          problemDraftModel.CreatedAt,
		UpdatedAt:          problemDraftModel.UpdatedAt,
	}

	for i, detail := range problemDraftModel.Details {
		response.Details[i] = ResponseProblemDraftDetail{
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
		response.Examples[i] = ResponseProblemDraftExample{
			Input:  example.Input,
			Output: example.Output,
		}
	}

	for i, displayName := range problemDifficulty.DisplayNames {
		response.ProblemDifficulty.DisplayNames[i] = ResponseProblemDifficultyDisplayName{
			Language: displayName.Language,
			Name:     displayName.DisplayName,
		}
	}

	return &response, nil
}
