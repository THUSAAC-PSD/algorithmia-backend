package database

import (
	"github.com/google/uuid"
)

type ProblemChatMessageAttachment struct {
	AttachmentID uuid.UUID `gorm:"primaryKey"`
	MessageID    uuid.UUID `gorm:"type:uuid"`
	MediaID      uuid.UUID `gorm:"type:uuid"`
}
