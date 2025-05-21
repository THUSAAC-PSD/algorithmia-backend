package login

import "github.com/google/uuid"

type User struct {
	UserID         uuid.UUID
	Username       string
	HashedPassword string
	Email          string
	IsSuperAdmin   bool
	Roles          []string
	Permissions    []string
}
