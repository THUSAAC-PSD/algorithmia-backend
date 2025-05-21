package contract

import (
	"context"

	"github.com/google/uuid"
)

type AuthUser struct {
	UserID uuid.UUID
	Email  string
}

type AuthUserDetails struct {
	Username     string
	IsSuperAdmin bool
	Roles        []string
	Permissions  []string
}

type AuthProvider interface {
	GetUser(ctx context.Context) (*AuthUser, error)

	Can(ctx context.Context, permissionNames ...string) (bool, error)

	MustGetUser(ctx context.Context) (AuthUser, error)
	MustGetUserDetails(ctx context.Context, userID uuid.UUID) (*AuthUserDetails, error)
}
