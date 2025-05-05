package createcontest

import (
	"time"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

var (
	ErrProblemCountRangeFlipped = errors.New("minProblemCount cannot be greater than maxProblemCount")
	ErrInvalidProblemCount      = errors.New("problem count must be greater than 1")
	ErrInvalidDeadlineDatetime  = errors.New("deadline datetime must be in the future")
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

func NewContest(
	title, description string,
	minProblemCount, maxProblemCount uint,
	deadlineDatetime time.Time,
) (Contest, error) {
	contestID, err := uuid.NewV7()
	if err != nil {
		return Contest{}, err
	}

	if minProblemCount > maxProblemCount {
		return Contest{}, errors.WithStack(ErrProblemCountRangeFlipped)
	}

	if minProblemCount < 1 || maxProblemCount < 1 {
		return Contest{}, errors.WithStack(ErrInvalidProblemCount)
	}

	now := time.Now()
	if deadlineDatetime.Before(now) {
		return Contest{}, errors.WithStack(ErrInvalidDeadlineDatetime)
	}

	return Contest{
		ContestID:        contestID,
		Title:            title,
		Description:      description,
		MinProblemCount:  minProblemCount,
		MaxProblemCount:  maxProblemCount,
		DeadlineDatetime: deadlineDatetime,
		CreatedAt:        now,
	}, nil
}
