package assigntesters

import "github.com/google/uuid"

type Command struct {
	ProblemID uuid.UUID   `param:"problem_id" validate:"required"`
	TesterIDs []uuid.UUID `                   validate:"required" json:"tester_ids"`
}
