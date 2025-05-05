package database

import "github.com/google/uuid"

type ProblemDifficultyDisplayName struct {
	DisplayNameID       uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProblemDifficultyID uuid.UUID `gorm:"type:uuid"`
	DisplayName         string
	Language            string
}
