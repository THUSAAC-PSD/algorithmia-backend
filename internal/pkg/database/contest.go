package database

import (
	"time"

	"github.com/google/uuid"
)

type Contest struct {
	ContestID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Title            string
	Description      string
	MinProblemCount  uint
	MaxProblemCount  uint
	DeadlineDatetime time.Time
	CreatedAt        time.Time
}
