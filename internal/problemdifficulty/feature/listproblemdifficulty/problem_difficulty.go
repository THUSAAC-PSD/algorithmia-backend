package listproblemdifficulty

import (
	"github.com/google/uuid"
)

type DisplayName struct {
	Language string `json:"language"`
	Name     string `json:"display_name"`
}

type ProblemDifficulty struct {
	ProblemDifficultyID uuid.UUID     `json:"problem_difficulty_id"`
	DisplayNames        []DisplayName `json:"display_names"`
}
