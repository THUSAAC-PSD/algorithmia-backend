package contract

import (
	"context"

	"github.com/google/uuid"
)

type AuthUser struct {
	UserID       uuid.UUID
	Email        string
	Username     string
	IsSuperAdmin bool
	Roles        []string
	Permissions  []string
}

type AuthProvider interface {
	GetUser(ctx context.Context) (*AuthUser, error)
	Can(ctx context.Context, permissionName string) (bool, error)

	MustGetUser(ctx context.Context) (AuthUser, error)
}
