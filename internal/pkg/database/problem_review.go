package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemReview struct {
	ProblemReviewID uuid.UUID `gorm:"primaryKey;type:uuid"`
	VersionID       uuid.UUID `gorm:"type:uuid;unique"`
	ReviewerID      uuid.UUID `gorm:"type:uuid"`
	Reviewer        User      `gorm:"foreignKey:ReviewerID"`
	Decision        string
	Comment         string
	CreatedAt       time.Time
}
