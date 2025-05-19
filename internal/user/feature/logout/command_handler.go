package logout

import (
	"context"

	"emperror.dev/errors"
)

type SessionManager interface {
	Delete(ctx context.Context) error
}

type CommandHandler struct {
	sessionManager SessionManager
}

func NewCommandHandler(sessionManager SessionManager) *CommandHandler {
	return &CommandHandler{
		sessionManager: sessionManager,
	}
}

func (c *CommandHandler) Handle(
	ctx context.Context,
) error {
	if err := c.sessionManager.Delete(ctx); err != nil {
		return errors.WrapIf(err, "failed to delete session")
	}

	return nil
}
