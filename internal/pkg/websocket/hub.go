package websocket

import (
	"context"
	"sync"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Hub struct {
	l logger.Logger

	r *Router

	register   chan *Client
	unregister chan *Client

	mu      sync.RWMutex // Protects clients and rooms maps
	clients map[*Client]bool
	rooms   map[uuid.UUID]map[*Client]bool // problemID -> set of clients focused on this problem's chat
}

func NewHub(l logger.Logger, r *Router) *Hub {
	return &Hub{
		l:          l,
		r:          r,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[uuid.UUID]map[*Client]bool),
	}
}

func (h *Hub) Run(ctx context.Context) {
	h.l.Info("WS Hub: Starting...")
	defer h.l.Info("WS Hub: Stopped.")

	for {
		select {
		case <-ctx.Done(): // Main application context is done
			h.shutdown()
			return
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			h.l.Infow("WS Hub: Client registered", map[string]interface{}{
				"user_id":       client.userID,
				"total_clients": len(h.clients),
			})
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send) // Signal client's writePump to stop

				oldProblemID := client.getFocusedProblemID()
				if oldProblemID.Valid {
					if clients, ok := h.rooms[oldProblemID.UUID]; ok {
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.rooms, oldProblemID.UUID)
						}

						h.l.Infow("WS Hub: Client unregistered from room", map[string]interface{}{
							"user_id":       client.userID,
							"problem_id":    oldProblemID.UUID,
							"total_clients": len(h.clients),
						})
					}
				}

				client.clearFocusedProblem()
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) shutdown() {
	h.l.Info("WS Hub: Shutting down all client connections...")

	h.mu.Lock() // Ensure no new registrations/unregistrations during shutdown
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)                                                       // Signal writePump
		_ = client.conn.Close(websocket.StatusGoingAway, "Server shutting down") // Politely close
	}

	h.clients = make(map[*Client]bool)
	h.rooms = make(map[uuid.UUID]map[*Client]bool)
}

func (h *Hub) processIncomingMessage(ctx context.Context, client *Client, rawMsg IncomingMessageBase) {
	switch rawMsg.Action {
	case "set-active-problem-chat":
		var payload SetActiveProblemChatClientPayload
		if err := h.r.UnmarshalPayload(rawMsg.Payload, &payload); err != nil {
			client.sendError("Invalid set-active-problem-chat payload", "INVALID_PAYLOAD", rawMsg.RequestID)
			return
		}

		// TODO: Authorization check here or in an application service
		// Can client.userID access/view payload.ProblemID?
		// For example:
		// canAccess, err := h.problemAuthService.CanUserAccessProblemChat(ctx, client.userID, payload.ProblemID)
		// if err != nil || !canAccess {
		//     client.sendError("Access denied to this problem chat.", "ACCESS_DENIED", rawMsg.RequestID)
		//     return
		// }

		newProblemID := payload.ProblemID
		oldProblemID := client.getFocusedProblemID()

		if oldProblemID.UUID == newProblemID {
			client.sendAck("Already focused on problem chat "+newProblemID.String(), rawMsg.RequestID)
			return
		}

		h.mu.Lock()
		if oldProblemID.Valid {
			if room, ok := h.rooms[oldProblemID.UUID]; ok {
				delete(room, client)
				if len(room) == 0 {
					delete(h.rooms, oldProblemID.UUID)
				}

				h.l.Infow("WS Hub: Client removed from room", map[string]interface{}{
					"user_id":    client.userID,
					"problem_id": oldProblemID.UUID,
				})
			}
		}

		if _, ok := h.rooms[newProblemID]; !ok {
			h.rooms[newProblemID] = make(map[*Client]bool)
		}

		h.rooms[newProblemID][client] = true
		h.mu.Unlock()

		client.setFocusedProblemID(newProblemID)
		client.sendAck("Active problem chat set to "+newProblemID.String(), rawMsg.RequestID)

		h.l.Infow("WS Hub: Client added to room", map[string]interface{}{
			"user_id":    client.userID,
			"problem_id": newProblemID,
		})

	default:
		if ok := h.r.Trigger(ctx, rawMsg.Action, rawMsg.Payload, func(message string, code ErrorCode) {
			client.sendError(message, code, rawMsg.RequestID)
		}, func(message string) {
			client.sendAck(message, rawMsg.RequestID)
		}); !ok {
			h.l.Infow("WS Hub: Unknown action", map[string]interface{}{
				"user_id": client.userID,
				"action":  rawMsg.Action,
			})

			client.sendError("Unknown action", ErrCodeUnknownAction, rawMsg.RequestID)
		}
	}
}

func (h *Hub) BroadcastToProblemChat(problemID uuid.UUID, messageBytes []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	roomClients, ok := h.rooms[problemID]
	if !ok || len(roomClients) == 0 {
		h.l.Infow("WS Hub Broadcast: No clients in problem chat room", map[string]interface{}{
			"problem_id": problemID,
		})

		return
	}

	h.l.Infow("WS Hub: Broadcasting message", map[string]interface{}{
		"problem_id":   problemID,
		"message":      string(messageBytes),
		"client_count": len(roomClients),
	})

	for client := range roomClients {
		client.sendRaw(messageBytes)
	}
}
