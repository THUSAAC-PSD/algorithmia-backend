package markcomplete

import "github.com/google/uuid"

type Command struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
}
