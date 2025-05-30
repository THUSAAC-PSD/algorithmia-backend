package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Username          string    `gorm:"unique"`
	Email             string    `gorm:"unique"`
	HashedPassword    string
	ProblemDrafts     []ProblemDraft       `gorm:"foreignKey:CreatorID"`
	Problems          []Problem            `gorm:"foreignKey:CreatorID"`
	Reviews           []ProblemReview      `gorm:"foreignKey:ReviewerID"`
	ProblemsReviewing []Problem            `gorm:"foreignKey:ReviewerID"`
	ChatMessages      []ProblemChatMessage `gorm:"foreignKey:SenderID"`
	Roles             []Role               `gorm:"many2many:user_roles"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
