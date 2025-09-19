package database

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationCode struct {
	EmailVerificationCodeID uuid.UUID `gorm:"primaryKey;type:uuid"`
	Code                    string
	Email                   string
	Username                string    // Store username for registration
	PasswordHash            string    // Store password hash for registration
	ExpiresAt               time.Time // Add expiration time
	CreatedAt               time.Time
}
