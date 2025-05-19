package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ProblemID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatorID         uuid.UUID `gorm:"type:uuid"`
	Status            string
	ProblemDraftID    uuid.UUID            `gorm:"type:uuid;unique"`
	TargetContestID   uuid.NullUUID        `gorm:"type:uuid"`
	AssignedContestID uuid.NullUUID        `gorm:"type:uuid"`
	ReviewerID        uuid.NullUUID        `gorm:"type:uuid"`
	TesterID          uuid.NullUUID        `gorm:"type:uuid"`
	ProblemVersions   []ProblemVersion     `gorm:"foreignKey:ProblemID"`
	ChatMessages      []ProblemChatMessage `gorm:"foreignKey:ProblemID"`
	CompletedAt       sql.NullTime
	CompletedBy       uuid.NullUUID `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}
