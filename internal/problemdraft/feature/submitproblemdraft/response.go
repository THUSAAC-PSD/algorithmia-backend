package submitproblemdraft

import (
	"github.com/google/uuid"
)

type Response struct {
	ProblemID        uuid.UUID `json:"problem_id"`
	ProblemVersionID uuid.UUID `json:"problem_version_id"`
}
