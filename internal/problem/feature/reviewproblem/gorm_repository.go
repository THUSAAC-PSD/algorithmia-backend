package reviewproblem

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
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

func (r *GormRepository) GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var pv database.ProblemVersion
	if err := db.WithContext(ctx).
		Model(&database.ProblemVersion{}).
		Select("problem_version_id").
		Where("problem_id = ?", problemID).
		Order("created_at DESC").
		First(&pv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, ErrProblemNotFound
		}

		return uuid.Nil, err
	}

	return pv.ProblemVersionID, nil
}

func (r *GormRepository) CreateReview(
	ctx context.Context,
	command *Command,
	reviewerID uuid.UUID,
	versionID uuid.UUID,
	createdAt time.Time,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	reviewID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to generate review ID")
	}

	if err := db.WithContext(ctx).
		Create(&database.ProblemReview{
			ProblemReviewID: reviewID,
			ReviewerID:      reviewerID,
			VersionID:       versionID,
			Decision:        string(command.Decision),
			Comment:         command.Comment,
			CreatedAt:       createdAt,
		}).Error; err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to create problem review")
	}

	return reviewID, nil
}

func (r *GormRepository) GetProblem(ctx context.Context, problemID uuid.UUID) (Problem, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problem Problem
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Select("status", "problem_draft_id AS draft_id").
		Where("problem_id = ?", problemID).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Problem{}, ErrProblemNotFound
		}

		return Problem{}, errors.WrapIf(err, "failed to get problem status")
	}

	return problem, nil
}

func (r *GormRepository) UpdateProblemStatus(
	ctx context.Context,
	problemID uuid.UUID,
	status constant.ProblemStatus,
) error {
	db := database.GetDBFromContext(ctx, r.db)

	if res := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ?", problemID).
		Update("status", status); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to update problem status")
	} else if res.RowsAffected == 0 {
		return errors.WithStack(ErrProblemNotFound)
	}

	return nil
}

func (r *GormRepository) SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	if res := db.WithContext(ctx).
		Model(&database.ProblemDraft{}).
		Where("problem_draft_id = ?", problemDraftID).
		Update("is_active", true); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to set problem draft active")
	} else if res.RowsAffected == 0 {
		return errors.New("failed to find problem draft")
	}

	return nil
}
