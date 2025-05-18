package infrastructure

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/reviewproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/testproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared/dto"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemActionRepository interface {
	GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error)
	CreateTestResult(
		ctx context.Context,
		command *testproblem.Command,
		testerID uuid.UUID,
		versionID uuid.UUID,
		createdAt time.Time,
	) (uuid.UUID, error)
	CreateReview(
		ctx context.Context,
		command *reviewproblem.Command,
		reviewerID uuid.UUID,
		versionID uuid.UUID,
		createdAt time.Time,
	) (uuid.UUID, error)
	GetProblem(ctx context.Context, problemID uuid.UUID) (dto.ProblemStatusAndVersion, error)
	UpdateProblemStatus(ctx context.Context, problemID uuid.UUID, status constant.ProblemStatus) error
	SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error
}

type ProblemActionGormRepository struct {
	db *gorm.DB
}

func NewProblemActionGormRepository(db *gorm.DB) *ProblemActionGormRepository {
	return &ProblemActionGormRepository{
		db: db,
	}
}

func (r *ProblemActionGormRepository) GetLatestProblemVersionID(
	ctx context.Context,
	problemID uuid.UUID,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var pv database.ProblemVersion
	if err := db.WithContext(ctx).
		Model(&database.ProblemVersion{}).
		Select("problem_version_id").
		Where("problem_id = ?", problemID).
		Order("created_at DESC").
		First(&pv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, shared.ErrProblemNotFound
		}

		return uuid.Nil, err
	}

	return pv.ProblemVersionID, nil
}

func (r *ProblemActionGormRepository) CreateTestResult(
	ctx context.Context,
	command *testproblem.Command,
	testerID uuid.UUID,
	versionID uuid.UUID,
	createdAt time.Time,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	resultID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to generate test result ID")
	}

	if err := db.WithContext(ctx).
		Create(&database.ProblemTestResult{
			ProblemTestResultID: resultID,
			TesterID:            testerID,
			VersionID:           versionID,
			Status:              string(command.Status),
			Comment:             command.Comment,
			CreatedAt:           createdAt,
		}).Error; err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to create problem test result")
	}

	return resultID, nil
}

func (r *ProblemActionGormRepository) CreateReview(
	ctx context.Context,
	command *reviewproblem.Command,
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

func (r *ProblemActionGormRepository) GetProblem(
	ctx context.Context,
	problemID uuid.UUID,
) (dto.ProblemStatusAndVersion, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problem dto.ProblemStatusAndVersion
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Select("status", "problem_draft_id AS draft_id").
		Where("problem_id = ?", problemID).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProblemStatusAndVersion{}, shared.ErrProblemNotFound
		}

		return dto.ProblemStatusAndVersion{}, errors.WrapIf(err, "failed to get problem status")
	}

	return problem, nil
}

func (r *ProblemActionGormRepository) UpdateProblemStatus(
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
		return errors.WithStack(shared.ErrProblemNotFound)
	}

	return nil
}

func (r *ProblemActionGormRepository) UpdateProblemReviewer(
	ctx context.Context,
	problemID uuid.UUID,
	reviewerID uuid.UUID,
) error {
	db := database.GetDBFromContext(ctx, r.db)

	if res := db.WithContext(ctx).
		Model(&database.Problem{}).
		Where("problem_id = ?", problemID).
		Update("reviewer_id", reviewerID); res.Error != nil {
		return errors.WrapIf(res.Error, "failed to update problem reviewer ID")
	} else if res.RowsAffected == 0 {
		return errors.WithStack(shared.ErrProblemNotFound)
	}

	return nil
}

func (r *ProblemActionGormRepository) SetProblemDraftActive(ctx context.Context, problemDraftID uuid.UUID) error {
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
