package createcontest

import "github.com/google/uuid"

type Response struct {
	ContestID uuid.UUID `json:"contest_id"`
}
