package submitproblemdraft

import "github.com/google/uuid"

type Command struct {
	ProblemDraftID  uuid.UUID     `param:"problem_draft_id" validate:"required"`
	TargetContestID uuid.NullUUID `                                             json:"target_contest_id"`
}
