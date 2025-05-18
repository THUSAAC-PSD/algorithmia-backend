package database

import (
	"time"

	"github.com/google/uuid"
)

type ProblemChatMessage struct {
	MessageID   uuid.UUID                      `gorm:"primaryKey"`
	ProblemID   uuid.UUID                      `gorm:"type:uuid"`
	SenderID    uuid.UUID                      `gorm:"type:uuid"`
	Content     string                         `gorm:"type:text"`
	Attachments []ProblemChatMessageAttachment `gorm:"foreignKey:MessageID"`
	CreatedAt   time.Time
}
