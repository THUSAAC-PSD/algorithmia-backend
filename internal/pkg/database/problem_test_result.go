package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemTestResult struct {
	ProblemTestResultID uuid.UUID `gorm:"primaryKey;type:uuid"`
	VersionID           uuid.UUID `gorm:"type:uuid;unique"`
	TesterID            uuid.UUID `gorm:"type:uuid"`
	Tester              User      `gorm:"foreignKey:TesterID"`
	Status              string
	Comment             string
	CreatedAt           time.Time
}
