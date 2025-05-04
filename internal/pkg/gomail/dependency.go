package gomail

import (
	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddGomail(container *dig.Container) error {
	err := container.Provide(func() (*Options, error) {
		return ProvideConfig()
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide gomail options")
	}

	return nil
}
