package database

import "time"

type EmailVerificationCode struct {
	EmailVerificationCodeID string `gorm:"primaryKey;type:uuid"`
	Code                    string
	Email                   string
	CreatedAt               time.Time
}
