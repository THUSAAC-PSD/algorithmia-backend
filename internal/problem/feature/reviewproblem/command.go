package reviewproblem

import "github.com/google/uuid"

type Decision string

const (
	DecisionApprove       Decision = "approve"
	DecisionReject        Decision = "reject"
	DecisionNeedsRevision Decision = "needs_revision"
)

type Command struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
	Decision  Decision  `                   validate:"required,oneof=approve reject needs_revision" json:"decision"`
	Comment   string    `                                                                           json:"comment"`
}
