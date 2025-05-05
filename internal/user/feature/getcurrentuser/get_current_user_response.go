package getcurrentuser

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type Response struct {
	User User `json:"user"`
}
