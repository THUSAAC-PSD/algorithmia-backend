package database

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	PermissionID uuid.UUID `gorm:"type:uuid;primary_key"`
	Name         string    `gorm:"unique"`
	Description  string
	CreatedAt    time.Time
	Roles        *[]Role `gorm:"many2many:role_permissions"`
}
