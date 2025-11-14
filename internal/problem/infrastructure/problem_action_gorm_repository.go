package infrastructure

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/dto"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/reviewproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/testproblem"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemActionRepository interface {
	GetLatestProblemVersionID(ctx context.Context, problemID uuid.UUID) (uuid.UUID, error)
	SaveTestResult(
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
	GetProblemTesterIDs(ctx context.Context, problemID uuid.UUID) ([]uuid.UUID, error)
	GetTestResultsForVersion(ctx context.Context, versionID uuid.UUID) ([]testproblem.ResultSummary, error)
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
			return uuid.Nil, problem.ErrProblemNotFound
		}

		return uuid.Nil, err
	}

	return pv.ProblemVersionID, nil
}

func (r *ProblemActionGormRepository) SaveTestResult(
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

	var existing database.ProblemTestResult
	err = db.WithContext(ctx).
		Where("tester_id = ? AND version_id = ?", testerID, versionID).
		First(&existing).Error
	if err == nil {
		existing.Status = string(command.Status)
		existing.Comment = command.Comment
		existing.CreatedAt = createdAt

		if err := db.WithContext(ctx).Save(&existing).Error; err != nil {
			return uuid.Nil, errors.WrapIf(err, "failed to update problem test result")
		}

		return existing.ProblemTestResultID, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, errors.WrapIf(err, "failed to get existing problem test result")
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

	var p dto.ProblemStatusAndVersion
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Select("status", "problem_draft_id AS draft_id").
		Where("problem_id = ?", problemID).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ProblemStatusAndVersion{}, problem.ErrProblemNotFound
		}

		return dto.ProblemStatusAndVersion{}, errors.WrapIf(err, "failed to get problem status")
	}

	return p, nil
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
		return errors.WithStack(problem.ErrProblemNotFound)
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
		return errors.WithStack(problem.ErrProblemNotFound)
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

func (r *ProblemActionGormRepository) GetProblemTesterIDs(ctx context.Context, problemID uuid.UUID) ([]uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	type testerRow struct {
		TesterID uuid.UUID
	}

	var rows []testerRow
	if err := db.WithContext(ctx).
		Table("problem_testers").
		Select("user_user_id AS tester_id").
		Where("problem_problem_id = ?", problemID).
		Scan(&rows).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get problem testers")
	}

	testerIDs := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		testerIDs = append(testerIDs, row.TesterID)
	}

	return testerIDs, nil
}

func (r *ProblemActionGormRepository) GetTestResultsForVersion(
	ctx context.Context,
	versionID uuid.UUID,
) ([]testproblem.ResultSummary, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var rows []database.ProblemTestResult
	if err := db.WithContext(ctx).
		Model(&database.ProblemTestResult{}).
		Select("tester_id", "status").
		Where("version_id = ?", versionID).
		Find(&rows).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get problem test results")
	}

	summaries := make([]testproblem.ResultSummary, len(rows))
	for i, row := range rows {
		summaries[i] = testproblem.ResultSummary{
			TesterID: row.TesterID,
			Status:   row.Status,
		}
	}

	return summaries, nil
}
