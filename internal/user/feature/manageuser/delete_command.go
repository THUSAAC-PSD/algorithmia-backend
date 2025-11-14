package manageuser

import "github.com/google/uuid"

type DeleteCommand struct {
	UserID uuid.UUID `param:"user_id" validate:"required"`
}
