package database

import "github.com/google/uuid"

type ProblemDifficulty struct {
	ProblemDifficultyID uuid.UUID                      `gorm:"primaryKey;type:uuid"`
	DisplayNames        []ProblemDifficultyDisplayName `gorm:"foreignKey:ProblemDifficultyID;constraint:OnDelete:CASCADE"`
	ProblemDrafts       []ProblemDraft                 `gorm:"foreignKey:ProblemDifficultyID"`
}
