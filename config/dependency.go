package config

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config/environment"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddAppConfig(container *dig.Container) error {
	err := container.Provide(func(environment environment.Environment) (*Config, error) {
		return NewAppConfig(environment)
	})

	return errors.WrapIf(err, "failed to provide app config")
}
