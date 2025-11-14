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
) error {
	db := database.GetDBFromContext(ctx, r.db)

	if !command.ProblemDraftID.Valid {
		return errors.WithStack(ErrInvalidProblemDraftID)
	}

	problemDraftModel := database.ProblemDraft{
		ProblemDraftID:      command.ProblemDraftID.UUID,
		ProblemDifficultyID: command.ProblemDifficultyID,
		CreatorID:           creatorID,
		Examples:            make([]database.ProblemDraftExample, len(command.Examples)),
		Details:             make([]database.ProblemDraftDetail, len(command.Details)),
		UpdatedAt:           updatedAt,
		IsActive:            true,
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

	if err := db.Transaction(func(tx *gorm.DB) error {
		if command.ProblemDifficultyID.Valid {
			if err := db.WithContext(ctx).
				Model(&database.ProblemDifficulty{}).
				Preload("DisplayNames").
				Where("problem_difficulty_id = ?", command.ProblemDifficultyID.UUID).
				First(&problemDifficulty).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.WithStack(ErrInvalidProblemDifficultyID)
				}

				return errors.WrapIf(err, "failed to get problem difficulty")
			}
		}

		if err := tx.WithContext(ctx).Where("problem_draft_id = ?", command.ProblemDraftID.UUID).Delete(&database.ProblemDraftExample{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete old problem draft examples")
		}

		if err := tx.WithContext(ctx).Where("problem_draft_id = ?", command.ProblemDraftID.UUID).Delete(&database.ProblemDraftDetail{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete old problem draft details")
		}

		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "problem_draft_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"problem_difficulty_id", "updated_at"}),
		}).Create(&problemDraftModel).Error; err != nil {
			return errors.WrapIf(err, "failed to upsert problem draft")
		}

		return nil
	}); err != nil {
		return errors.WrapIf(err, "failed to upsert problem draft in transaction")
	}

	return nil
}

func (r *GormRepository) VerifyActiveProblemDraftCreator(
	ctx context.Context,
	problemDraftID uuid.UUID,
	creatorID uuid.UUID,
) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.ProblemDraft{}).
		Where("problem_draft_id = ? AND creator_id = ?", problemDraftID, creatorID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to count problem drafts")
	}

	return count > 0, nil
}
