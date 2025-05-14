package dto

import (
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

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
	IsActive           bool                  `json:"is_active"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

func FromGormProblemDifficulty(problemDifficulty database.ProblemDifficulty) ProblemDifficulty {
	dto := ProblemDifficulty{
		ProblemDifficultyID: problemDifficulty.ProblemDifficultyID,
		DisplayNames:        make([]ProblemDifficultyDisplayName, len(problemDifficulty.DisplayNames)),
	}

	for i, displayName := range problemDifficulty.DisplayNames {
		dto.DisplayNames[i] = ProblemDifficultyDisplayName{
			Language: displayName.Language,
			Name:     displayName.DisplayName,
		}
	}

	return dto
}

func FromGormProblemDraft(problemDraft database.ProblemDraft, problemDifficulty ProblemDifficulty) ProblemDraft {
	dto := ProblemDraft{
		ProblemDraftID:    problemDraft.ProblemDraftID,
		ProblemDifficulty: problemDifficulty,
		CreatorID:         problemDraft.CreatorID,
		Details:           make([]ProblemDraftDetail, len(problemDraft.Details)),
		Examples:          make([]ProblemDraftExample, len(problemDraft.Examples)),
		IsActive:          problemDraft.IsActive,
		CreatedAt:         problemDraft.CreatedAt,
		UpdatedAt:         problemDraft.UpdatedAt,
	}

	if problemDraft.SubmittedProblem.ProblemID != uuid.Nil {
		dto.SubmittedProblemID = uuid.NullUUID{Valid: true, UUID: problemDraft.SubmittedProblem.ProblemID}
	}

	for i, detail := range problemDraft.Details {
		dto.Details[i] = ProblemDraftDetail{
			Language:     detail.Language,
			Title:        detail.Title,
			Background:   detail.Background,
			Statement:    detail.Statement,
			InputFormat:  detail.InputFormat,
			OutputFormat: detail.OutputFormat,
			Note:         detail.Note,
		}
	}

	for i, example := range problemDraft.Examples {
		dto.Examples[i] = ProblemDraftExample{
			Input:  example.Input,
			Output: example.Output,
		}
	}

	return dto
}
