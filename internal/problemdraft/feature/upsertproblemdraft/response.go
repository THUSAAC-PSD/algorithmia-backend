package upsertproblemdraft

import (
	"time"

	"github.com/google/uuid"
)

type ResponseProblemDraftDetail struct {
	Language     string `json:"language"`
	Title        string `json:"title"`
	Background   string `json:"background"`
	Statement    string `json:"statement"`
	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Note         string `json:"note"`
}

type ResponseProblemDraftExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ResponseProblemDifficultyDisplayName struct {
	Language string `json:"language"`
	Name     string `json:"display_name"`
}

type ResponseProblemDifficulty struct {
	ProblemDifficultyID uuid.UUID                              `json:"problem_difficulty_id"`
	DisplayNames        []ResponseProblemDifficultyDisplayName `json:"display_names"`
}

type ResponseProblemDraft struct {
	ProblemDraftID     uuid.UUID                     `json:"problem_draft_id"`
	ProblemDifficulty  ResponseProblemDifficulty     `json:"problem_difficulty"`
	CreatorID          uuid.UUID                     `json:"creator_id"`
	Details            []ResponseProblemDraftDetail  `json:"details"`
	Examples           []ResponseProblemDraftExample `json:"examples"`
	SubmittedProblemID uuid.NullUUID                 `json:"submitted_problem_id"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
}

type Response struct {
	ProblemDraft ResponseProblemDraft `json:"problem_draft"`
}
