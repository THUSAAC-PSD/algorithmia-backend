package assigntester

import "github.com/google/uuid"

type Command struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
	UserID    uuid.UUID `                   validate:"required" json:"user_id"`
}
