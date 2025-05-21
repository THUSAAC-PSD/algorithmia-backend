package getproblem

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"
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

func (r *GormRepository) GetProblem(
	ctx context.Context,
	problemID uuid.UUID,
) (*ResponseProblem, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var p database.Problem
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Preload("Creator").
		Preload("Reviewer").
		Preload("Testers").
		Preload("ProblemVersions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("ProblemVersions.Details").
		Preload("ProblemVersions.Examples").
		Preload("ProblemVersions.Review").
		Preload("ProblemVersions.Review.Reviewer").
		Preload("ProblemVersions.TestResult").
		Preload("ProblemVersions.TestResult.Tester").
		Preload("ProblemVersions.ProblemDifficulty").
		Preload("ProblemVersions.ProblemDifficulty.DisplayNames").
		Preload("TargetContest").
		Preload("AssignedContest").
		Where("problem_id = ?", problemID).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(problem.ErrProblemNotFound)
		}

		return nil, errors.WrapIf(err, "failed to get problem")
	}

	result := &ResponseProblem{
		ProblemID:       p.ProblemID,
		LatestVersionID: p.ProblemVersions[0].ProblemVersionID,
		Versions:        make([]ResponseProblemVersion, 0, len(p.ProblemVersions)),
		Status:          constant.FromStringToProblemStatus(p.Status),
		Creator: ResponseUser{
			UserID:   p.Creator.UserID,
			Username: p.Creator.Username,
		},
		Testers:   make([]ResponseUser, 0, len(p.Testers)),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	for _, version := range p.ProblemVersions {
		v := ResponseProblemVersion{
			VersionID:         version.ProblemVersionID,
			ProblemDifficulty: dto.FromGormProblemDifficulty(version.ProblemDifficulty),
			Details:           make([]ResponseProblemDetail, 0, len(version.Details)),
			Examples:          make([]ResponseProblemExample, 0, len(version.Examples)),
			CreatedAt:         version.CreatedAt,
		}

		if version.Review != nil {
			v.Review = &ResponseReview{
				ReviewerID: version.Review.Reviewer.UserID,
				Comment:    version.Review.Comment,
				Decision:   version.Review.Decision,
				CreatedAt:  version.Review.CreatedAt,
			}
		}

		if version.TestResult != nil {
			v.TestResult = &ResponseTestResult{
				TesterID:  version.TestResult.Tester.UserID,
				Comment:   version.TestResult.Comment,
				Status:    version.TestResult.Status,
				CreatedAt: version.TestResult.CreatedAt,
			}
		}

		for _, detail := range version.Details {
			v.Details = append(v.Details, ResponseProblemDetail{
				Language:     detail.Language,
				Title:        detail.Title,
				Background:   detail.Background,
				Statement:    detail.Statement,
				InputFormat:  detail.InputFormat,
				OutputFormat: detail.OutputFormat,
				Note:         detail.Note,
			})
		}

		for _, example := range version.Examples {
			v.Examples = append(v.Examples, ResponseProblemExample{
				Input:  example.Input,
				Output: example.Output,
			})
		}

		result.Versions = append(result.Versions, v)
	}

	for _, tester := range p.Testers {
		result.Testers = append(result.Testers, ResponseUser{
			UserID:   tester.UserID,
			Username: tester.Username,
		})
	}

	if p.Reviewer != nil {
		result.Reviewer = &ResponseUser{
			UserID:   p.Reviewer.UserID,
			Username: p.Reviewer.Username,
		}
	}

	if p.TargetContest != nil {
		result.TargetContest = &ResponseContest{
			ContestID: p.TargetContest.ContestID,
			Title:     p.TargetContest.Title,
		}
	}

	if p.AssignedContest != nil {
		result.AssignedContest = &ResponseContest{
			ContestID: p.AssignedContest.ContestID,
			Title:     p.AssignedContest.Title,
		}
	}

	return result, nil
}
