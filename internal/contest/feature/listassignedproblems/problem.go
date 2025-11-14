package listassignedproblems

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ProblemID         uuid.UUID            `json:"problem_id"`
	Title             []ProblemDetailTitle `json:"title" gorm:"-"`
	ProblemDifficulty ProblemDifficulty    `json:"problem_difficulty" gorm:"-"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

type ProblemDetailTitle struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type ProblemDifficulty struct {
	ProblemDifficultyID uuid.UUID                      `json:"problem_difficulty_id"`
	DisplayNames        []ProblemDifficultyDisplayName `json:"display_names" gorm:"-"`
}

type ProblemDifficultyDisplayName struct {
	Language    string `json:"language"`
	DisplayName string `json:"display_name"`
}
