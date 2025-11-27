package manageuser

import "github.com/google/uuid"

type ResetPasswordCommand struct {
	UserID          uuid.UUID `param:"user_id" validate:"required"`
	NewPassword     string    `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string    `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
