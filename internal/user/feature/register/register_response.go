package register

import (
	"time"

	"github.com/google/uuid"
)

type ResponseUser struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	User ResponseUser `json:"user"`
}
