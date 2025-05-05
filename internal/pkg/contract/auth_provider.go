package contract

import (
	"context"

	"github.com/google/uuid"
)

type AuthUser struct {
	UserID   uuid.UUID
	Email    string
	Username string
}

type AuthProvider interface {
	GetUser(ctx context.Context) (*AuthUser, error)

	MustGetUser(ctx context.Context) (AuthUser, error)
}
