package reviewproblem

import (
	"github.com/google/uuid"
)

type Response struct {
	ReviewID         uuid.UUID `json:"review_id"`
	ProblemVersionID uuid.UUID `json:"problem_version_id"`
}
