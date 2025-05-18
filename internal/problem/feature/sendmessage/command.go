package sendmessage

import "github.com/google/uuid"

type Command struct {
	ProblemID          uuid.UUID   `validate:"required" json:"problem_id"`
	Content            string      `validate:"required" json:"content"`
	AttachmentMediaIDs []uuid.UUID `validate:"required" json:"attachment_media_ids"`
}
