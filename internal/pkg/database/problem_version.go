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
	Details             []ProblemVersionDetail  `gorm:"foreignKey:ProblemVersionID"`
	Examples            []ProblemVersionExample `gorm:"foreignKey:ProblemVersionID"`
	CreatedAt           time.Time
}
