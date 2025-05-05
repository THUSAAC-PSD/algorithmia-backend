package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemDraft struct {
	ProblemDraftID      uuid.UUID             `gorm:"primaryKey;type:uuid"`
	CreatorID           uuid.UUID             `gorm:"type:uuid"`
	ProblemDifficultyID uuid.NullUUID         `gorm:"type:uuid"`
	Examples            []ProblemDraftExample `gorm:"foreignKey:ProblemDraftID"`
	Details             []ProblemDraftDetail  `gorm:"foreignKey:ProblemDraftID"`
	IsActive            bool                  // True until the problem draft is submitted, then false if it needs revision
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time

	// TODO: link to Problems
}
