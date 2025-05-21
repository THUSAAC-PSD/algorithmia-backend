package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemDraft struct {
	ProblemDraftID      uuid.UUID     `gorm:"primaryKey;type:uuid"`
	CreatorID           uuid.UUID     `gorm:"type:uuid"`
	ProblemDifficultyID uuid.NullUUID `gorm:"type:uuid"`
	ProblemDifficulty   ProblemDifficulty
	SubmittedProblem    Problem               `gorm:"foreignKey:ProblemDraftID"`
	Examples            []ProblemDraftExample `gorm:"foreignKey:ProblemDraftID"`
	Details             []ProblemDraftDetail  `gorm:"foreignKey:ProblemDraftID"`
	IsActive            bool                  // False after the problem draft is submitted, then true again when it needs revision
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Deleted             gorm.DeletedAt `gorm:"index"`
}
