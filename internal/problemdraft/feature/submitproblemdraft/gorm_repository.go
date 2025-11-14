package submitproblemdraft

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/dto"

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

func (r *GormRepository) GetProblemDraft(ctx context.Context, problemDraftID uuid.UUID) (*dto.ProblemDraft, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problemDraft database.ProblemDraft
	if err := db.WithContext(ctx).
		Preload("Examples").
		Preload("Details").
		Preload("ProblemDifficulty").
		Preload("ProblemDifficulty.DisplayNames").
		Preload("SubmittedProblem").
		Where("problem_draft_id = ?", problemDraftID).
		First(&problemDraft).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrProblemDraftNotFound)
		}

		return nil, errors.WrapIf(err, "failed to get problem draft")
	}

	response := dto.FromGormProblemDraft(problemDraft, dto.FromGormProblemDifficulty(problemDraft.ProblemDifficulty))
	return &response, nil
}

func (r *GormRepository) GetProblemStatus(ctx context.Context, problemID uuid.UUID) (*constant.ProblemStatus, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problem database.Problem
	if err := db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Select("status").
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.WrapIf(err, "failed to get problem status")
	}

	status := constant.FromStringToProblemStatus(problem.Status)
	return &status, nil
}

func (r *GormRepository) SetProblemDraftInactive(ctx context.Context, problemDraftID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	if err := db.WithContext(ctx).
		Model(&database.ProblemDraft{}).
		Where("problem_draft_id = ?", problemDraftID).
		Update("is_active", false).Error; err != nil {
		return errors.WrapIf(err, "failed to set problem draft inactive")
	}

	return nil
}

func (r *GormRepository) UpsertProblemFromDraft(
	ctx context.Context,
	draft *dto.ProblemDraft,
	targetContestID uuid.NullUUID,
	status constant.ProblemStatus,
	updatedAt time.Time,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	problem := database.Problem{
		ProblemDraftID:  draft.ProblemDraftID,
		CreatorID:       draft.CreatorID,
		Status:          string(status),
		TargetContestID: targetContestID,
		CreatedAt:       updatedAt,
		UpdatedAt:       updatedAt,
	}

	if draft.SubmittedProblemID.Valid {
		problem.ProblemID = draft.SubmittedProblemID.UUID
	} else {
		problemID, err := uuid.NewV7()
		if err != nil {
			return uuid.Nil, errors.WrapIf(err, "failed to generate new problem ID")
		}

		problem.ProblemID = problemID
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if targetContestID.Valid {
			var contest database.Contest
			if err := tx.WithContext(ctx).
				Where("contest_id = ?", targetContestID.UUID).
				First(&contest).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.WithStack(ErrContestNotFound)
				}

				return errors.WrapIf(err, "failed to check contest existence")
			}
		}

		if err := tx.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "problem_draft_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"updated_at", "status", "target_contest_id"}),
			}).
			Create(&problem).Error; err != nil {
			return errors.WrapIf(err, "failed to upsert problem")
		}

		return nil
	}); err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to run transaction")
	}

	return problem.ProblemID, nil
}

func (r *GormRepository) CreateProblemVersionFromDraft(
	ctx context.Context,
	problemID uuid.UUID,
	draft *dto.ProblemDraft,
	createdAt time.Time,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	problemVersionID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to generate new problem version ID")
	}

	problemVersion := database.ProblemVersion{
		ProblemVersionID:    problemVersionID,
		ProblemID:           problemID,
		ProblemDifficultyID: draft.ProblemDifficulty.ProblemDifficultyID,
		SubmittedBy:         draft.CreatorID,
		Details:             make([]database.ProblemVersionDetail, len(draft.Details)),
		Examples:            make([]database.ProblemVersionExample, len(draft.Examples)),
		CreatedAt:           createdAt,
	}

	for i, detail := range draft.Details {
		detailID, err := uuid.NewV7()
		if err != nil {
			return uuid.Nil, errors.WrapIf(err, "failed to generate new problem version detail ID")
		}

		problemVersion.Details[i] = database.ProblemVersionDetail{
			DetailID:         detailID,
			ProblemVersionID: problemVersionID,
			Language:         detail.Language,
			Title:            detail.Title,
			Background:       detail.Background,
			Statement:        detail.Statement,
			InputFormat:      detail.InputFormat,
			OutputFormat:     detail.OutputFormat,
			Note:             detail.Note,
		}
	}

	for i, example := range draft.Examples {
		exampleID, err := uuid.NewV7()
		if err != nil {
			return uuid.Nil, errors.WrapIf(err, "failed to generate new problem version example ID")
		}

		problemVersion.Examples[i] = database.ProblemVersionExample{
			ExampleID:        exampleID,
			ProblemVersionID: problemVersionID,
			Input:            example.Input,
			Output:           example.Output,
		}
	}

	if err := db.WithContext(ctx).
		Create(&problemVersion).Error; err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to create problem version")
	}

	return problemVersion.ProblemVersionID, nil
}
