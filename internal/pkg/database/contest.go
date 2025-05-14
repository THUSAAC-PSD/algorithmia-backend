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
	TargetedProblems []Problem `gorm:"foreignKey:TargetContestID"`
	AssignedProblems []Problem `gorm:"foreignKey:AssignedContestID"`
	DeadlineDatetime time.Time
	CreatedAt        time.Time
	DeletedAt        time.Time
}
