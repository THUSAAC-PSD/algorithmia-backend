package logger

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddLogger(container *dig.Container) error {
	err := container.Provide(func(environment environment.Environment) (*Options, error) {
		return ProvideConfig(environment)
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide log options")
	}

	err = container.Provide(func(opts *Options, environment environment.Environment) Logger {
		return NewZapLogger(opts, environment)
	})
	return errors.WrapIf(err, "failed to provide logger")
}
