package deletecontest

import "github.com/google/uuid"

type Command struct {
	ContestID uuid.UUID `param:"contest_id" validate:"required"`
}
