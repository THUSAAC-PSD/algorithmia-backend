package logout

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/mehdihadeli/go-mediatr"
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
	command *Command,
) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := c.sessionManager.Delete(ctx); err != nil {
		return mediatr.Unit{}, errors.WithStack(err)
	}

	return mediatr.Unit{}, nil
}
