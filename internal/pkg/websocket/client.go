package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

const (
	maxMessageSize = 2048

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	l      logger.Logger
	hub    *Hub
	conn   *websocket.Conn
	userID uuid.UUID

	// Buffered channel of outbound messages.
	// Messages placed here will be picked up by the writePump.
	send chan []byte

	mu               sync.RWMutex
	focusedProblemID uuid.NullUUID
}

func NewClient(l logger.Logger, hub *Hub, conn *websocket.Conn, userID uuid.UUID) *Client {
	return &Client{
		l:      l,
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close(websocket.StatusNormalClosure, "Client read pump ended")

		c.l.Infow("WS Client disconnected", map[string]interface{}{
			"user_id": c.userID,
		})
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		var rawMessage IncomingMessageBase

		if err := wsjson.Read(ctx, c.conn, &rawMessage); err != nil {
			closeStatus := websocket.CloseStatus(err)
			if closeStatus == websocket.StatusNormalClosure || closeStatus == websocket.StatusGoingAway ||
				errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, io.EOF) {
				c.l.Infow("WS Client read error", map[string]interface{}{
					"user_id": c.userID,
					"error":   err,
				})
			} else {
				c.l.Errorw("WS Client read error", map[string]interface{}{
					"user_id": c.userID,
					"error":   err,
				})
			}
			break
		}

		msgCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		c.hub.processIncomingMessage(msgCtx, c, rawMessage)

		cancel()
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.l.Infow("WS Client write pump ended", map[string]interface{}{
			"user_id": c.userID,
		})
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.Close(websocket.StatusNormalClosure, "Hub closed send channel")
				c.l.Infow("WS Client write pump: Hub closed send channel", map[string]interface{}{
					"user_id": c.userID,
				})

				return
			}

			wCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(wCtx, websocket.MessageText, message)

			cancel()

			if err != nil {
				c.l.Errorw("WS Client write error", map[string]interface{}{
					"user_id": c.userID,
				})

				return
			}
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait/2) // Shorter timeout for ping
			if err := c.conn.Ping(pingCtx); err != nil {
				cancel()

				c.l.Errorw("WS Client ping error", map[string]interface{}{
					"user_id": c.userID,
				})

				return
			}

			cancel()
		case <-ctx.Done():
			c.l.Infow("WS Client context done", map[string]interface{}{
				"user_id": c.userID,
			})

			_ = c.conn.Close(websocket.StatusGoingAway, "Server shutting down or client context cancelled")
			return
		}
	}
}

func (c *Client) sendRaw(message []byte) {
	select {
	case c.send <- message:
	default: // Don't block if client's send buffer is full
		c.l.Infow("WS Client send buffer full", map[string]interface{}{
			"user_id": c.userID,
		})
	}
}

func (c *Client) sendError(errMsg string, code ErrorCode, requestID string) {
	env := OutgoingMessageEnvelope{
		Type: contract.MessageTypeError,
		Payload: ErrorServerPayload{
			Message: errMsg,
			Code:    code,
		},
		RequestID: requestID,
	}

	bytes, err := json.Marshal(env)
	if err != nil {
		c.l.Errorw("WS Client %s: Failed to marshal error for client", map[string]interface{}{
			"user_id": c.userID,
			"error":   err,
		})

		return
	}

	c.sendRaw(bytes)
}

func (c *Client) sendAck(message string, requestID string) {
	env := OutgoingMessageEnvelope{
		Type: contract.MessageTypeAck,
		Payload: AckServerPayload{
			Message: message,
		},
		RequestID: requestID,
	}
	bytes, err := json.Marshal(env)
	if err != nil {
		c.l.Errorw("WS Client %s: Failed to marshal ack for client", map[string]interface{}{
			"user_id": c.userID,
			"error":   err,
		})

		return
	}

	c.sendRaw(bytes)
}

func (c *Client) setFocusedProblemID(problemID uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.focusedProblemID = uuid.NullUUID{Valid: true, UUID: problemID}
}

func (c *Client) getFocusedProblemID() uuid.NullUUID {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.focusedProblemID
}

func (c *Client) clearFocusedProblem() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.focusedProblemID = uuid.NullUUID{Valid: false}
}
