package testproblem

import (
	"github.com/google/uuid"
)

type Response struct {
	TestResultID     uuid.UUID `json:"test_result_id"`
	ProblemVersionID uuid.UUID `json:"problem_version_id"`
}
