package listassignedproblems

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type assignedProblemRecord struct {
	ProblemID           uuid.UUID     `gorm:"column:problem_id"`
	CreatedAt           time.Time     `gorm:"column:created_at"`
	UpdatedAt           time.Time     `gorm:"column:updated_at"`
	LatestVersionID     uuid.NullUUID `gorm:"column:latest_version_id"`
	ProblemDifficultyID uuid.NullUUID `gorm:"column:problem_difficulty_id"`
}

type problemVersionDetailRecord struct {
	ProblemVersionID uuid.UUID `gorm:"column:problem_version_id"`
	Language         string    `gorm:"column:language"`
	Title            string    `gorm:"column:title"`
}

type problemDifficultyDisplayRecord struct {
	ProblemDifficultyID uuid.UUID `gorm:"column:problem_difficulty_id"`
	Language            string    `gorm:"column:language"`
	DisplayName         string    `gorm:"column:display_name"`
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (g *GormRepository) GetAssignedProblems(ctx context.Context, contestID uuid.UUID) ([]Problem, error) {
	var records []assignedProblemRecord
	query := g.db.WithContext(ctx).
		Table("problems p").
		Joins(`
			LEFT JOIN LATERAL (
				SELECT
					pv.problem_version_id,
					pv.problem_difficulty_id,
					pv.created_at
				FROM problem_versions pv
				WHERE pv.problem_id = p.problem_id
				ORDER BY pv.created_at DESC
				LIMIT 1
			) latest ON TRUE
		`).
		Select(`
			p.problem_id,
			p.created_at,
			p.updated_at,
			latest.problem_version_id AS latest_version_id,
			latest.problem_difficulty_id AS problem_difficulty_id
		`).
		Where("p.assigned_contest_id = ?", contestID).
		Order("p.created_at ASC")

	if err := query.Scan(&records).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get assigned problems")
	}

	if len(records) == 0 {
		return []Problem{}, nil
	}

	versionIDSet := make(map[uuid.UUID]struct{})
	difficultyIDSet := make(map[uuid.UUID]struct{})
	for _, rec := range records {
		if rec.LatestVersionID.Valid {
			versionIDSet[rec.LatestVersionID.UUID] = struct{}{}
		}
		if rec.ProblemDifficultyID.Valid {
			difficultyIDSet[rec.ProblemDifficultyID.UUID] = struct{}{}
		}
	}

	versionIDs := make([]uuid.UUID, 0, len(versionIDSet))
	for id := range versionIDSet {
		versionIDs = append(versionIDs, id)
	}

	difficultyIDs := make([]uuid.UUID, 0, len(difficultyIDSet))
	for id := range difficultyIDSet {
		difficultyIDs = append(difficultyIDs, id)
	}

	titlesByVersionID := make(map[uuid.UUID][]ProblemDetailTitle)
	if len(versionIDs) > 0 {
		var detailRows []problemVersionDetailRecord
		if err := g.db.WithContext(ctx).
			Table("problem_version_details").
			Select("problem_version_id, language, title").
			Where("problem_version_id IN ?", versionIDs).
			Find(&detailRows).Error; err != nil {
			return nil, errors.Wrap(err, "failed to get problem titles")
		}

		for _, row := range detailRows {
			titlesByVersionID[row.ProblemVersionID] = append(titlesByVersionID[row.ProblemVersionID], ProblemDetailTitle{
				Language: row.Language,
				Title:    row.Title,
			})
		}
	}

	displayNamesByDifficultyID := make(map[uuid.UUID][]ProblemDifficultyDisplayName)
	if len(difficultyIDs) > 0 {
		var difficultyRows []problemDifficultyDisplayRecord
		if err := g.db.WithContext(ctx).
			Table("problem_difficulty_display_names").
			Select("problem_difficulty_id, language, display_name").
			Where("problem_difficulty_id IN ?", difficultyIDs).
			Find(&difficultyRows).Error; err != nil {
			return nil, errors.Wrap(err, "failed to get difficulty display names")
		}

		for _, row := range difficultyRows {
			displayNamesByDifficultyID[row.ProblemDifficultyID] = append(displayNamesByDifficultyID[row.ProblemDifficultyID], ProblemDifficultyDisplayName{
				Language:    row.Language,
				DisplayName: row.DisplayName,
			})
		}
	}

	problems := make([]Problem, 0, len(records))
	for _, rec := range records {
		titles := titlesByVersionID[rec.LatestVersionID.UUID]
		if titles == nil {
			titles = make([]ProblemDetailTitle, 0)
		}

		displayNames := displayNamesByDifficultyID[rec.ProblemDifficultyID.UUID]
		if displayNames == nil {
			displayNames = make([]ProblemDifficultyDisplayName, 0)
		}

		problem := Problem{
			ProblemID: rec.ProblemID,
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
			Title:     titles,
			ProblemDifficulty: ProblemDifficulty{
				ProblemDifficultyID: rec.ProblemDifficultyID.UUID,
				DisplayNames:        displayNames,
			},
		}

		problems = append(problems, problem)
	}

	return problems, nil
}
