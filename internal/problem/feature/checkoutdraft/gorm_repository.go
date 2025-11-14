package checkoutdraft

import (
	"context"
	"time"

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
	return &GormRepository{db: db}
}

func (r *GormRepository) GetProblemSummary(ctx context.Context, problemID uuid.UUID) (ProblemSummary, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problem database.Problem
	if err := db.WithContext(ctx).
		Select("problem_id", "creator_id", "problem_draft_id").
		Where("problem_id = ?", problemID).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProblemSummary{}, errors.WithStack(ErrProblemNotFound)
		}
		return ProblemSummary{}, errors.WrapIf(err, "failed to get problem")
	}

	return ProblemSummary{
		ProblemID:      problem.ProblemID,
		ProblemDraftID: problem.ProblemDraftID,
		CreatorID:      problem.CreatorID,
	}, nil
}

func (r *GormRepository) GetLatestVersion(ctx context.Context, problemID uuid.UUID) (*VersionAggregate, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var version database.ProblemVersion
	if err := db.WithContext(ctx).
		Preload("Details").
		Preload("Examples").
		Where("problem_id = ?", problemID).
		Order("created_at DESC").
		First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrProblemNotFound)
		}
		return nil, errors.WrapIf(err, "failed to get latest problem version")
	}

	v := &VersionAggregate{
		ProblemVersionID:    version.ProblemVersionID,
		ProblemDifficultyID: version.ProblemDifficultyID,
		Details:             make([]VersionDetail, len(version.Details)),
		Examples:            make([]VersionExample, len(version.Examples)),
	}

	for i, detail := range version.Details {
		v.Details[i] = VersionDetail{
			Language:     detail.Language,
			Title:        detail.Title,
			Background:   detail.Background,
			Statement:    detail.Statement,
			InputFormat:  detail.InputFormat,
			OutputFormat: detail.OutputFormat,
			Note:         detail.Note,
		}
	}

	for i, example := range version.Examples {
		v.Examples[i] = VersionExample{
			Input:  example.Input,
			Output: example.Output,
		}
	}

	return v, nil
}

func (r *GormRepository) ReplaceDraftFromVersion(
	ctx context.Context,
	draftID uuid.UUID,
	version *VersionAggregate,
	updatedAt time.Time,
) error {
	db := database.GetDBFromContext(ctx, r.db)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("problem_draft_id = ?", draftID).Delete(&database.ProblemDraftDetail{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete draft details")
		}

		if err := tx.Where("problem_draft_id = ?", draftID).Delete(&database.ProblemDraftExample{}).Error; err != nil {
			return errors.WrapIf(err, "failed to delete draft examples")
		}

		details := make([]database.ProblemDraftDetail, len(version.Details))
		for i, detail := range version.Details {
			details[i] = database.ProblemDraftDetail{
				DetailID:       uuid.Must(uuid.NewV7()),
				ProblemDraftID: draftID,
				Language:       detail.Language,
				Title:          detail.Title,
				Background:     detail.Background,
				Statement:      detail.Statement,
				InputFormat:    detail.InputFormat,
				OutputFormat:   detail.OutputFormat,
				Note:           detail.Note,
			}
		}

		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return errors.WrapIf(err, "failed to insert draft details")
			}
		}

		examples := make([]database.ProblemDraftExample, len(version.Examples))
		for i, example := range version.Examples {
			examples[i] = database.ProblemDraftExample{
				ExampleID:      uuid.Must(uuid.NewV7()),
				ProblemDraftID: draftID,
				Input:          example.Input,
				Output:         example.Output,
			}
		}

		if len(examples) > 0 {
			if err := tx.Create(&examples).Error; err != nil {
				return errors.WrapIf(err, "failed to insert draft examples")
			}
		}

		update := map[string]interface{}{
			"is_active":  true,
			"updated_at": updatedAt,
		}

		if version.ProblemDifficultyID != uuid.Nil {
			update["problem_difficulty_id"] = uuid.NullUUID{
				UUID:  version.ProblemDifficultyID,
				Valid: true,
			}
		} else {
			update["problem_difficulty_id"] = uuid.NullUUID{}
		}

		if err := tx.Model(&database.ProblemDraft{}).
			Where("problem_draft_id = ?", draftID).
			Updates(update).Error; err != nil {
			return errors.WrapIf(err, "failed to update problem draft")
		}

		return nil
	})
}

func (r *GormRepository) GetProblemDraft(ctx context.Context, draftID uuid.UUID) (*dto.ProblemDraft, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var draft database.ProblemDraft
	if err := db.WithContext(ctx).
		Preload("ProblemDifficulty").
		Preload("ProblemDifficulty.DisplayNames").
		Preload("Details").
		Preload("Examples").
		Preload("SubmittedProblem").
		Where("problem_draft_id = ?", draftID).
		First(&draft).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get problem draft")
	}

	dtoDraft := dto.FromGormProblemDraft(
		draft,
		dto.FromGormProblemDifficulty(draft.ProblemDifficulty),
	)

	return &dtoDraft, nil
}
