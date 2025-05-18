package database

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	MediaID   uuid.UUID `gorm:"primaryKey"`
	URL       string
	FileName  string
	MIMEType  string
	FileSize  uint64
	CreatedAt time.Time
}
