package upsertproblemdraft

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared/dto"
)

type Response struct {
	ProblemDraft dto.ProblemDraft `json:"problem_draft"`
}
