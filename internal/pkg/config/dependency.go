package config

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config/environment"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddEnv(container *dig.Container, environments ...environment.Environment) error {
	err := container.Provide(func() environment.Environment {
		return environment.ConfigEnv(environments...)
	})
	return errors.WrapIf(err, "failed to provide environment")
}
