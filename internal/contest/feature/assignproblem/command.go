package assignproblem

import "github.com/google/uuid"

type Command struct {
	ContestID uuid.UUID `param:"contest_id" validate:"required"`
	ProblemID uuid.UUID `                   validate:"required" json:"problem_id"`
}
