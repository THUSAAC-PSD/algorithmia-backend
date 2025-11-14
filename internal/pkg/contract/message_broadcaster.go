package contract

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeAck   MessageType = "ack"
	MessageTypeError MessageType = "error"

	MessageTypeUser      MessageType = "user"
	MessageTypeSubmitted MessageType = "submitted"
	MessageTypeEdited    MessageType = "edited"
	MessageTypeReviewed  MessageType = "reviewed"
	MessageTypeTested    MessageType = "tested"
	MessageTypeCompleted MessageType = "completed"
)

type MessageUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type MessageAttachment struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	MIMEType string `json:"mime_type"`
	Size     uint64 `json:"size"`
}

type MessageBroadcaster interface {
	BroadcastUserMessage(
		problemID uuid.UUID,
		messageID uuid.UUID,
		content string,
		sender MessageUser,
		attachments []MessageAttachment,
		timestamp time.Time,
	) error

	BroadcastSubmittedMessage(
		problemID uuid.UUID,
		submitter MessageUser,
		timestamp time.Time,
	) error

	BroadcastEditedMessage(
		problemID uuid.UUID,
		editor MessageUser,
		timestamp time.Time,
	) error

	BroadcastReviewedMessage(
		problemID uuid.UUID,
		reviewer MessageUser,
		decision string,
		timestamp time.Time,
	) error

	BroadcastTestedMessage(
		problemID uuid.UUID,
		tester MessageUser,
		status string,
		timestamp time.Time,
	) error

	BroadcastCompletedMessage(
		problemID uuid.UUID,
		completer MessageUser,
		timestamp time.Time,
	) error
}
