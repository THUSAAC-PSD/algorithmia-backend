package database

import "github.com/google/uuid"

type ProblemVersionExample struct {
	ExampleID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProblemVersionID uuid.UUID `gorm:"type:uuid"`
	Input            string
	Output           string
}
