package logger

import (
	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddLogger(container *dig.Container) error {
	if err := container.Provide(NewZapLogger); err != nil {
		return errors.WrapIf(err, "failed to provide logger")
	}

	return nil
}
