package getproblem

import "github.com/google/uuid"

type Query struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
}
