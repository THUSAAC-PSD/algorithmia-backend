package database

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationCode struct {
	EmailVerificationCodeID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Code                    string
	Email                   string
	CreatedAt               time.Time
}
