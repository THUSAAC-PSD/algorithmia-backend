package testproblem

import "github.com/google/uuid"

type Status string

const (
	StatusPassed Status = "passed"
	StatusFailed Status = "failed"
)

type Command struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
	Status    Status    `                   validate:"required,oneof=passed failed" json:"status"`
	Comment   string    `                                                           json:"comment"`
}
