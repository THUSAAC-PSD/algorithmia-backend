package listmessage

import (
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"github.com/google/uuid"
)

type ResponseChatMessage struct {
	MessageType string      `json:"message_type"`
	Payload     interface{} `json:"payload"`
	Timestamp   time.Time   `json:"timestamp"`
}

type ResponseChatUserPayload struct {
	MessageID   uuid.UUID                    `json:"message_id"`
	Sender      contract.MessageUser         `json:"sender"`
	Content     string                       `json:"content"`
	Attachments []contract.MessageAttachment `json:"attachments"`
}

type ResponseChatSubmittedPayload struct {
	Submitter contract.MessageUser `json:"submitter"`
	// ChangedFields is present for edited messages to summarize what changed
	ChangedFields []string `json:"changed_fields,omitempty"`
	// Diffs contains per-field before/after content for edited messages
	Diffs map[string]FieldDiff `json:"diffs,omitempty"`
}

// FieldDiff represents a simple textual before/after for a field
type FieldDiff struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type ResponseChatReviewedPayload struct {
	Reviewer contract.MessageUser `json:"reviewer"`
	Decision string               `json:"decision"`
}

type ResponseChatTestedPayload struct {
	Tester contract.MessageUser `json:"tester"`
	Status string               `json:"status"`
}

type ResponseChatCompletedPayload struct {
	Completer contract.MessageUser `json:"completer"`
	Status    string               `json:"status"`
}

type Response struct {
	Messages []ResponseChatMessage `json:"messages"`
}
