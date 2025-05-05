package listproblemdraft

import "github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared/dto"

type Response struct {
	// TODO: Only return summary information
	ProblemDrafts []dto.ProblemDraft `json:"problem_drafts"`
}
