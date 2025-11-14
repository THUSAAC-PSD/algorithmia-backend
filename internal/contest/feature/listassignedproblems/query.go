package listassignedproblems

import "github.com/google/uuid"

type Query struct {
	ContestID uuid.UUID `param:"contest_id" validate:"required"`
}
