package resetpassword

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
}

type User struct {
	UserID         uuid.UUID
	HashedPassword string
}
