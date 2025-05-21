package database

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string
	Description  string
	IsSuperAdmin bool
	Permissions  *[]Permission `gorm:"many2many:role_permissions"`
	Users        *[]User       `gorm:"many2many:user_roles"`
	CreatedAt    time.Time
}
