package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	Deleted          gorm.DeletedAt `gorm:"index"`
}
