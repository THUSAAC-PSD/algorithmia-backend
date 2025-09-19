package listproblem

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
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

type flatProblemData struct {
	ProblemID        uuid.UUID `gorm:"column:problem_id"`
	ProblemStatus    string    `gorm:"column:problem_status"`
	ProblemCreatedAt time.Time `gorm:"column:problem_created_at"`
	ProblemUpdatedAt time.Time `gorm:"column:problem_updated_at"`
	ProblemDraftID   uuid.UUID `gorm:"column:problem_draft_id"`

	CreatorID       uuid.UUID `gorm:"column:creator_id"`
	CreatorUsername string    `gorm:"column:creator_username"`

	ReviewerID       uuid.NullUUID  `gorm:"column:reviewer_id"`
	ReviewerUsername sql.NullString `gorm:"column:reviewer_username"`

	TargetContestID    uuid.NullUUID  `gorm:"column:target_contest_id"`
	TargetContestTitle sql.NullString `gorm:"column:target_contest_title"`

	AssignedContestID    uuid.NullUUID  `gorm:"column:assigned_contest_id"`
	AssignedContestTitle sql.NullString `gorm:"column:assigned_contest_title"`

	LatestVersionID           uuid.NullUUID `gorm:"column:latest_version_id"`
	LatestVersionDifficultyID uuid.NullUUID `gorm:"column:latest_version_difficulty_id"`
}

type flatProblemTesterData struct {
	ProblemID  uuid.UUID `gorm:"column:problem_id"`
	TesterID   uuid.UUID `gorm:"column:tester_id"`
	TesterName string    `gorm:"column:tester_username"`
}

func (r *GormRepository) GetAllRelatedProblems(
	ctx context.Context,
	userID uuid.UUID,
	showAll bool,
	showCreated bool,
	showAllPendingReview bool,
	showAssignedTesting bool,
	onlyShowCompleted bool,
) ([]ResponseProblem, error) {
	db := database.GetDBFromContext(ctx, r.db)

	flatProblems, err := r.fetchRelatedProblems(
		ctx,
		db,
		userID,
		showAll,
		showCreated,
		showAllPendingReview,
		showAssignedTesting,
		onlyShowCompleted,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch related problems")
	}

	var problemTesters []flatProblemTesterData
	if err := db.WithContext(ctx).
		Table("problem_testers pt").
		Joins("LEFT JOIN users ON users.user_id = pt.user_user_id").
		Select("pt.problem_problem_id as problem_id, pt.user_user_id as tester_id, users.username as tester_username").
		Scan(&problemTesters).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch problem testers")
	}

	testerByProblemID := make(map[uuid.UUID][]ResponseUser)
	for _, p := range flatProblems {
		testerByProblemID[p.ProblemID] = make([]ResponseUser, 0)
	}

	for _, pt := range problemTesters {
		testerByProblemID[pt.ProblemID] = append(testerByProblemID[pt.ProblemID], ResponseUser{
			UserID:   pt.TesterID,
			Username: pt.TesterName,
		})
	}

	latestVersionIDs := make([]uuid.UUID, 0, len(flatProblems))
	difficultyIDs := make([]uuid.UUID, 0, len(flatProblems))

	for _, problem := range flatProblems {
		if problem.LatestVersionID.Valid {
			latestVersionIDs = append(latestVersionIDs, problem.LatestVersionID.UUID)
		}

		if problem.LatestVersionDifficultyID.Valid {
			difficultyIDs = append(difficultyIDs, problem.LatestVersionDifficultyID.UUID)
		}
	}

	titlesByVersionID, err := r.fetchTitlesForVersions(ctx, db, latestVersionIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch titles for versions")
	}

	difficultiesByID, err := r.fetchDifficulties(ctx, db, difficultyIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch difficulties")
	}

	result := make([]ResponseProblem, 0, len(flatProblems))
	for _, problem := range flatProblems {
		p := ResponseProblem{
			ProblemID: problem.ProblemID,
			Status:    constant.FromStringToProblemStatus(problem.ProblemStatus),
			Creator: ResponseUser{
				UserID:   problem.CreatorID,
				Username: problem.CreatorUsername,
			},
			Testers:   testerByProblemID[problem.ProblemID],
			CreatedAt: problem.ProblemCreatedAt,
			UpdatedAt: problem.ProblemUpdatedAt,
			Titles:    make([]ResponseProblemTitle, 0),
		}

		if problem.ReviewerID.Valid && problem.ReviewerUsername.Valid {
			p.Reviewer = &ResponseUser{
				UserID:   problem.ReviewerID.UUID,
				Username: problem.ReviewerUsername.String,
			}
		}

		if problem.TargetContestID.Valid && problem.TargetContestTitle.Valid {
			p.TargetContest = &ResponseContest{
				ContestID: problem.TargetContestID.UUID,
				Title:     problem.TargetContestTitle.String,
			}
		}

		if problem.AssignedContestID.Valid && problem.AssignedContestTitle.Valid {
			p.AssignedContest = &ResponseContest{
				ContestID: problem.AssignedContestID.UUID,
				Title:     problem.AssignedContestTitle.String,
			}
		}

		if problem.LatestVersionID.Valid {
			p.Titles = titlesByVersionID[problem.LatestVersionID.UUID]
		}

		if problem.LatestVersionDifficultyID.Valid {
			p.ProblemDifficulty = difficultiesByID[problem.LatestVersionDifficultyID.UUID]
		}

		result = append(result, p)
	}

	return result, nil
}

func (r *GormRepository) fetchRelatedProblems(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
	showAll bool,
	showCreated bool,
	showAllPendingReview bool,
	showAssignedTesting bool,
	onlyShowCompleted bool,
) ([]flatProblemData, error) {
	var problems []flatProblemData
	query := db.WithContext(ctx).
		Table("problems p").
		Joins("LEFT JOIN users creator_u ON creator_u.user_id = p.creator_id").
		Joins("LEFT JOIN users reviewer_u ON reviewer_u.user_id = p.reviewer_id").
		Joins("LEFT JOIN contests target_c ON target_c.contest_id = p.target_contest_id").
		Joins("LEFT JOIN contests assigned_c ON assigned_c.contest_id = p.assigned_contest_id").
		Joins("LEFT JOIN problem_testers pr ON pr.problem_problem_id = p.problem_id").
		Joins(`
			LEFT JOIN LATERAL (
				SELECT
					pv.problem_version_id,
					pv.problem_difficulty_id,
					pv.created_at
				FROM
					problem_versions pv
				WHERE
					pv.problem_id = p.problem_id
				ORDER BY
					pv.created_at DESC
				LIMIT 1
			) rpv ON TRUE
		`).
		Select(`
            p.problem_id,
            p.status as problem_status,
            p.created_at as problem_created_at,
            p.updated_at as problem_updated_at,
            p.problem_draft_id,
            p.creator_id,
            p.reviewer_id,
            p.target_contest_id,
            p.assigned_contest_id,
            creator_u.username as creator_username,
            reviewer_u.username as reviewer_username,
            target_c.title as target_contest_title,
            assigned_c.title as assigned_contest_title,
			rpv.problem_version_id as latest_version_id,
            rpv.problem_difficulty_id as latest_version_difficulty_id 
        `)

	// We build OR visibility predicates; user should ALWAYS see their own created problems
	visibilityPredicates := make([]string, 0)
	visibilityArgs := make([]interface{}, 0)

	// Always own problems
	visibilityPredicates = append(visibilityPredicates, "p.creator_id = ?")
	visibilityArgs = append(visibilityArgs, userID)

	// As reviewer
	visibilityPredicates = append(visibilityPredicates, "p.reviewer_id = ?")
	visibilityArgs = append(visibilityArgs, userID)

	// Testing assignments
	if showAssignedTesting {
		visibilityPredicates = append(visibilityPredicates, "pr.user_user_id = ?")
		visibilityArgs = append(visibilityArgs, userID)
	}

	// Pending review for reviewers with global view
	if showAllPendingReview {
		visibilityPredicates = append(visibilityPredicates, "p.status = ?")
		visibilityArgs = append(visibilityArgs, constant.ProblemStatusPendingReview)
	}

	// Global list all overrides everything else
	if showAll {
		visibilityPredicates = []string{"1 = 1"}
		visibilityArgs = []interface{}{}
	}

	whereClause := strings.Join(visibilityPredicates, " OR ")
	if onlyShowCompleted {
		whereClause = "(" + whereClause + ") AND p.status = ?"
		visibilityArgs = append(visibilityArgs, constant.ProblemStatusCompleted)
	}

	fmt.Printf(
		"fetchRelatedProblems | user=%s showAll=%v showCreated=%v showAllPendingReview=%v showAssignedTesting=%v onlyShowCompleted=%v where=%s args=%v\n",
		userID,
		showAll,
		showCreated,
		showAllPendingReview,
		showAssignedTesting,
		onlyShowCompleted,
		whereClause,
		visibilityArgs,
	)

	if err := query.Where(whereClause, visibilityArgs...).Scan(&problems).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all related problems")
	}

	return problems, nil
}

func (r *GormRepository) fetchTitlesForVersions(
	ctx context.Context,
	db *gorm.DB,
	versionIDs []uuid.UUID,
) (map[uuid.UUID][]ResponseProblemTitle, error) {
	var details []database.ProblemVersionDetail
	if err := db.WithContext(ctx).
		Model(&database.ProblemVersionDetail{}).
		Where("problem_version_id IN ?", versionIDs).
		Find(&details).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch problem version details")
	}

	titlesMap := make(map[uuid.UUID][]ResponseProblemTitle)
	for _, d := range details {
		titlesMap[d.ProblemVersionID] = append(titlesMap[d.ProblemVersionID], ResponseProblemTitle{
			Language: d.Language,
			Title:    d.Title,
		})
	}

	return titlesMap, nil
}

func (r *GormRepository) fetchDifficulties(
	ctx context.Context,
	db *gorm.DB,
	difficultyIDs []uuid.UUID,
) (map[uuid.UUID]dto.ProblemDifficulty, error) {
	var dbDifficulties []database.ProblemDifficulty
	if err := db.WithContext(ctx).
		Model(&database.ProblemDifficulty{}).
		Preload("DisplayNames").
		Where("problem_difficulty_id IN ?", difficultyIDs).
		Find(&dbDifficulties).Error; err != nil {
		return nil, errors.Wrap(err, "failed to fetch problem difficulties with display names")
	}

	difficultiesMap := make(map[uuid.UUID]dto.ProblemDifficulty)
	for _, dbDiff := range dbDifficulties {
		displayNamesDTO := make([]dto.ProblemDifficultyDisplayName, len(dbDiff.DisplayNames))
		for i, dn := range dbDiff.DisplayNames {
			displayNamesDTO[i] = dto.ProblemDifficultyDisplayName{
				Language: dn.Language,
				Name:     dn.DisplayName,
			}
		}

		difficultiesMap[dbDiff.ProblemDifficultyID] = dto.ProblemDifficulty{
			ProblemDifficultyID: dbDiff.ProblemDifficultyID,
			DisplayNames:        displayNamesDTO,
		}
	}

	return difficultiesMap, nil
}
