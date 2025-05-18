package websocket

import (
	"encoding/json"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/google/uuid"
)

type WsBroadcaster struct {
	hub *Hub
	l   logger.Logger
}

func NewWsBroadcaster(hub *Hub, l logger.Logger) *WsBroadcaster {
	return &WsBroadcaster{
		hub: hub,
		l:   l,
	}
}

func (b *WsBroadcaster) BroadcastUserMessage(
	problemID uuid.UUID,
	messageID uuid.UUID,
	content string,
	sender contract.MessageUser,
	attachments []contract.MessageAttachment,
	timestamp time.Time,
) error {
	b.l.Infow("WS Broadcaster: Broadcasting user message", map[string]interface{}{
		"problem_id": problemID,
		"sender_id":  sender.UserID,
	})

	payload := NewMessageServerPayload{
		ID:        messageID,
		ProblemID: problemID,
		Sender: UserServerPayload{
			ID:       sender.UserID,
			Username: sender.Username,
		},
		Content:     content,
		Attachments: make([]AttachmentServerInfo, 0, len(attachments)),
		Timestamp:   timestamp,
	}

	for _, attachment := range attachments {
		payload.Attachments = append(payload.Attachments, AttachmentServerInfo{
			URL:      attachment.URL,
			FileName: attachment.FileName,
			MIMEType: attachment.MIMEType,
			Size:     attachment.Size,
		})
	}

	envelope := OutgoingMessageEnvelope{
		Type:    contract.MessageTypeUser,
		Payload: payload,
	}

	if err := b.broadcastEnvelope(problemID, &envelope); err != nil {
		return errors.WrapIf(err, "failed to broadcast message")
	}
	return nil
}

func (b *WsBroadcaster) BroadcastSubmittedMessage(
	problemID uuid.UUID,
	submitter contract.MessageUser,
	timestamp time.Time,
) error {
	b.l.Infow("WS Broadcaster: Broadcasting submitted message", map[string]interface{}{
		"problem_id":   problemID,
		"submitter_id": submitter.UserID,
	})

	payload := SubmittedMessageServerPayload{
		ProblemID: problemID,
		Submitter: UserServerPayload{
			ID:       submitter.UserID,
			Username: submitter.Username,
		},
		Timestamp: timestamp,
	}

	envelope := OutgoingMessageEnvelope{
		Type:    contract.MessageTypeSubmitted,
		Payload: payload,
	}

	if err := b.broadcastEnvelope(problemID, &envelope); err != nil {
		return errors.WrapIf(err, "failed to broadcast message")
	}
	return nil
}

func (b *WsBroadcaster) BroadcastReviewedMessage(
	problemID uuid.UUID,
	reviewer contract.MessageUser,
	decision string,
	timestamp time.Time,
) error {
	b.l.Infow("WS Broadcaster: Broadcasting reviewed message", map[string]interface{}{
		"problem_id":  problemID,
		"reviewer_id": reviewer.UserID,
	})

	payload := ReviewedMessageServerPayload{
		ProblemID: problemID,
		Reviewer: UserServerPayload{
			ID:       reviewer.UserID,
			Username: reviewer.Username,
		},
		Decision:  decision,
		Timestamp: timestamp,
	}

	envelope := OutgoingMessageEnvelope{
		Type:    contract.MessageTypeReviewed,
		Payload: payload,
	}

	if err := b.broadcastEnvelope(problemID, &envelope); err != nil {
		return errors.WrapIf(err, "failed to broadcast message")
	}
	return nil
}

func (b *WsBroadcaster) BroadcastTestedMessage(
	problemID uuid.UUID,
	tester contract.MessageUser,
	status string,
	timestamp time.Time,
) error {
	b.l.Infow("WS Broadcaster: Broadcasting tested message", map[string]interface{}{
		"problem_id": problemID,
		"tester_id":  tester.UserID,
	})

	payload := TestedMessageServerPayload{
		ProblemID: problemID,
		Tester: UserServerPayload{
			ID:       tester.UserID,
			Username: tester.Username,
		},
		Status:    status,
		Timestamp: timestamp,
	}

	envelope := OutgoingMessageEnvelope{
		Type:    contract.MessageTypeTested,
		Payload: payload,
	}

	if err := b.broadcastEnvelope(problemID, &envelope); err != nil {
		return errors.WrapIf(err, "failed to broadcast message")
	}
	return nil
}

func (b *WsBroadcaster) BroadcastCompletedMessage(
	problemID uuid.UUID,
	completer contract.MessageUser,
	timestamp time.Time,
) error {
	b.l.Infow("WS Broadcaster: Broadcasting completed message", map[string]interface{}{
		"problem_id":   problemID,
		"completer_id": completer.UserID,
	})

	payload := CompletedMessageServerPayload{
		ProblemID: problemID,
		Completer: UserServerPayload{
			ID:       completer.UserID,
			Username: completer.Username,
		},
		Timestamp: timestamp,
	}

	envelope := OutgoingMessageEnvelope{
		Type:    contract.MessageTypeCompleted,
		Payload: payload,
	}

	if err := b.broadcastEnvelope(problemID, &envelope); err != nil {
		return errors.WrapIf(err, "failed to broadcast message")
	}
	return nil
}

func (b *WsBroadcaster) broadcastEnvelope(problemID uuid.UUID, e *OutgoingMessageEnvelope) error {
	messageBytes, err := json.Marshal(e)
	if err != nil {
		return errors.WrapIf(err, "failed to marshal message for broadcast")
	}

	b.hub.BroadcastToProblemChat(problemID, messageBytes)
	return nil
}
