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
	MessageTypeReviewed  MessageType = "reviewed"
	MessageTypeTested    MessageType = "tested"
	MessageTypeCompleted MessageType = "completed"
)

type MessageUser struct {
	UserID   uuid.UUID
	Username string
}

type MessageAttachment struct {
	URL      string
	FileName string
	MIMEType string
	Size     uint64
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
