package deleteproblemdraft

import "github.com/google/uuid"

type Command struct {
	ProblemDraftID uuid.UUID `param:"problem_draft_id" validate:"required"`
}
