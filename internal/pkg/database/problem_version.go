package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemVersion struct {
	ProblemVersionID    uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProblemID           uuid.UUID `gorm:"type:uuid"`
	ProblemDifficultyID uuid.UUID `gorm:"type:uuid"`
	ProblemDifficulty   ProblemDifficulty
	SubmittedBy         uuid.UUID               `gorm:"type:uuid"`
	SubmittedByUser     User                    `gorm:"foreignKey:SubmittedBy"`
	Details             []ProblemVersionDetail  `gorm:"foreignKey:ProblemVersionID"`
	Examples            []ProblemVersionExample `gorm:"foreignKey:ProblemVersionID"`
	Review              *ProblemReview          `gorm:"foreignKey:VersionID"`
	TestResults         []ProblemTestResult     `gorm:"foreignKey:VersionID"`
	CreatedAt           time.Time
}
