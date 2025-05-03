package register

import (
	"time"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID
	Username       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

func NewUser(username, email, hashedPassword string) (User, error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return User{}, errors.WrapIf(err, "failed to generate user ID")
	}

	return User{
		UserID:         userID,
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}, nil
}
