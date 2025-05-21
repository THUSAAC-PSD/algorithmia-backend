package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ProblemID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatorID         uuid.UUID `gorm:"type:uuid"`
	Creator           User      `gorm:"foreignKey:CreatorID"`
	Status            string
	ProblemDraftID    uuid.UUID            `gorm:"type:uuid;unique"`
	TargetContestID   uuid.NullUUID        `gorm:"type:uuid"`
	TargetContest     *Contest             `gorm:"foreignKey:TargetContestID"`
	AssignedContestID uuid.NullUUID        `gorm:"type:uuid"`
	AssignedContest   *Contest             `gorm:"foreignKey:AssignedContestID"`
	ReviewerID        uuid.NullUUID        `gorm:"type:uuid"`
	Reviewer          *User                `gorm:"foreignKey:ReviewerID"`
	Testers           []User               `gorm:"many2many:problem_testers"`
	ProblemVersions   []ProblemVersion     `gorm:"foreignKey:ProblemID"`
	ChatMessages      []ProblemChatMessage `gorm:"foreignKey:ProblemID"`
	CompletedAt       sql.NullTime
	CompletedBy       uuid.NullUUID `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
