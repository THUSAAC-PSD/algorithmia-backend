package websocket

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddWebsocket(container *dig.Container) error {
	if err := container.Provide(NewRouter); err != nil {
		return errors.WrapIf(err, "failed to provide websocket router")
	}

	if err := container.Provide(NewHub); err != nil {
		return errors.WrapIf(err, "failed to provide websocket hub")
	}

	if err := container.Provide(NewWsBroadcaster,
		dig.As(new(contract.MessageBroadcaster))); err != nil {
		return errors.WrapIf(err, "failed to provide websocket broadcaster")
	}

	return nil
}
