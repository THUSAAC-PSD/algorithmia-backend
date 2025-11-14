package manageuser

import "github.com/google/uuid"

type UpdateCommand struct {
	UserID   uuid.UUID `param:"user_id" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Roles    []string  `json:"roles" validate:"required,min=1,dive,required"`
}
