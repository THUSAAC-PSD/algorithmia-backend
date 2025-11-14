package manageuser

import (
	"time"

	"github.com/google/uuid"
)

type ResponseUser struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseRole struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsSuperAdmin bool   `json:"is_super_admin"`
}

type ListResponse struct {
	Users []ResponseUser `json:"users"`
	Roles []ResponseRole `json:"roles"`
}

type UpdateResponse struct {
	User ResponseUser `json:"user"`
}
