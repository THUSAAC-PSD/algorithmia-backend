package database

import "github.com/google/uuid"

type UserRole struct {
	UserRoleID uuid.UUID `gorm:"primaryKey;type:uuid"`
	RoleType   string
	UserID     uuid.UUID `gorm:"type:uuid"`
}
