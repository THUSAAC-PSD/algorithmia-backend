package listcontest

import (
	"time"

	"github.com/google/uuid"
)

type Contest struct {
	ContestID        uuid.UUID `json:"contest_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	MinProblemCount  uint      `json:"min_problem_count"`
	MaxProblemCount  uint      `json:"max_problem_count"`
	DeadlineDatetime time.Time `json:"deadline_datetime"`
	CreatedAt        time.Time `json:"created_at"`
}
