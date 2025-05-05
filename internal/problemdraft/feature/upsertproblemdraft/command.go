package upsertproblemdraft

import "github.com/google/uuid"

type CommandDetail struct {
	Language     string `json:"language"      validate:"required"`
	Title        string `json:"title"`
	Background   string `json:"background"`
	Statement    string `json:"statement"`
	InputFormat  string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Note         string `json:"note"`
}

type CommandExample struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Command struct {
	ProblemDraftID      uuid.NullUUID    `json:"problem_draft_id"`
	ProblemDifficultyID uuid.NullUUID    `json:"problem_difficulty_id"`
	Details             []CommandDetail  `json:"details"               validate:"required"`
	Examples            []CommandExample `json:"examples"`
}
