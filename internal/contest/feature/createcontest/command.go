package createcontest

import "time"

type Command struct {
	Title            string    `json:"title"             validate:"required"`
	Description      string    `json:"description"       validate:"required"`
	MinProblemCount  uint      `json:"min_problem_count" validate:"required,gte=1"`
	MaxProblemCount  uint      `json:"max_problem_count" validate:"required,gte=1"`
	DeadlineDatetime time.Time `json:"deadline_datetime" validate:"required"`
}
