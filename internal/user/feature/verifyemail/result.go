package verifyemail

import "github.com/google/uuid"

type Result struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
