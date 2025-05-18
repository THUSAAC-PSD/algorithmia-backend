package websocket

import (
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"github.com/google/uuid"
)

type IncomingMessageBase struct {
	Action    string      `json:"action"`
	Payload   interface{} `json:"payload"`
	RequestID string      `json:"request_id"`
}

type OutgoingMessageEnvelope struct {
	Type      contract.MessageType `json:"type"`
	Payload   interface{}          `json:"payload"`
	RequestID string               `json:"request_id,omitempty"`
}

// ErrorServerPayload for sending errors back to the client
type ErrorServerPayload struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

// AckServerPayload for general acknowledgements
type AckServerPayload struct {
	Message string `json:"message,omitempty"`
}

// NewMessageServerPayload for broadcasting new chat messages
type NewMessageServerPayload struct {
	ID          uuid.UUID              `json:"id"`
	ProblemID   uuid.UUID              `json:"problem_id"`
	Sender      UserServerPayload      `json:"sender"`
	Content     string                 `json:"content"`
	Attachments []AttachmentServerInfo `json:"attachments"`
	Timestamp   time.Time              `json:"timestamp"`
}

// SubmittedMessageServerPayload for problem submission messages
type SubmittedMessageServerPayload struct {
	ProblemID uuid.UUID         `json:"problem_id"`
	Submitter UserServerPayload `json:"submitter"`
	Timestamp time.Time         `json:"timestamp"`
}

// ReviewedMessageServerPayload for problem review messages
type ReviewedMessageServerPayload struct {
	ProblemID uuid.UUID         `json:"problem_id"`
	Reviewer  UserServerPayload `json:"reviewer"`
	Decision  string            `json:"decision"`
	Timestamp time.Time         `json:"timestamp"`
}

// TestedMessageServerPayload for problem testing messages
type TestedMessageServerPayload struct {
	ProblemID uuid.UUID         `json:"problem_id"`
	Tester    UserServerPayload `json:"tester"`
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
}

// CompletedMessageServerPayload for problem completion messages
type CompletedMessageServerPayload struct {
	ProblemID uuid.UUID         `json:"problem_id"`
	Completer UserServerPayload `json:"completer"`
	Timestamp time.Time         `json:"timestamp"`
}

// UserServerPayload for user information in messages
type UserServerPayload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// AttachmentServerInfo for detailed attachment information in messages
type AttachmentServerInfo struct {
	URL      string `json:"url"`
	FileName string `json:"filename"`
	MIMEType string `json:"mime_type"`
	Size     uint64 `json:"size"`
}

type SetActiveProblemChatClientPayload struct {
	ProblemID uuid.UUID `json:"problem_id"`
}
