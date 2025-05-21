package getcurrentuser

import "github.com/google/uuid"

type ResponseUser struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	Roles        []string  `json:"roles"`
	Permissions  []string  `json:"permissions"`
}

type Response struct {
	User ResponseUser `json:"user"`
}
