package getcurrentuser

import "github.com/google/uuid"

type ResponseUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Roles    []string  `json:"roles"`
}

type Response struct {
	User ResponseUser `json:"user"`
}
