package database

import "github.com/google/uuid"

type ProblemVersionDetail struct {
	DetailID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProblemVersionID uuid.UUID `gorm:"type:uuid"`
	Language         string
	Title            string
	Background       string
	Statement        string
	InputFormat      string
	OutputFormat     string
	Note             string
}
