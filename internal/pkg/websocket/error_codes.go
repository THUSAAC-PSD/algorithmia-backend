package websocket

type ErrorCode string

const (
	ErrCodeInvalidPayload      ErrorCode = "invalid_payload"
	ErrCodeInternalServerError ErrorCode = "internal_server_error"
	ErrCodeUnknownAction       ErrorCode = "unknown_action"
)
