package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProblemDraftDetail struct {
	Language     string `json:"language"`
	Title        string `json:"title"`
	Background   string `json:"background"`
	Statement    string `json:"statement"`
	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Note         string `json:"note"`
}

type ProblemDraftExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ProblemDifficultyDisplayName struct {
	Language string `json:"language"`
	Name     string `json:"display_name"`
}

type ProblemDifficulty struct {
	ProblemDifficultyID uuid.UUID                      `json:"problem_difficulty_id"`
	DisplayNames        []ProblemDifficultyDisplayName `json:"display_names"`
}

type ProblemDraft struct {
	ProblemDraftID     uuid.UUID             `json:"problem_draft_id"`
	ProblemDifficulty  ProblemDifficulty     `json:"problem_difficulty"`
	CreatorID          uuid.UUID             `json:"creator_id"`
	Details            []ProblemDraftDetail  `json:"details"`
	Examples           []ProblemDraftExample `json:"examples"`
	SubmittedProblemID uuid.NullUUID         `json:"submitted_problem_id"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}
