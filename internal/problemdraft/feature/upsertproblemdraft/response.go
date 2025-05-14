package upsertproblemdraft

import "github.com/google/uuid"

type Response struct {
	ProblemDraftID uuid.UUID `json:"problem_draft_id"`
}
