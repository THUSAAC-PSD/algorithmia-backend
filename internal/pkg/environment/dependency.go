package environment

import (
	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddEnv(container *dig.Container, environments ...Environment) error {
	err := container.Provide(func() Environment {
		return ConfigEnv(environments...)
	})
	return errors.WrapIf(err, "failed to provide environment")
}
