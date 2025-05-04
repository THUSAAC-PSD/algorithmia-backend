package database

import (
	"emperror.dev/errors"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func AddGorm(container *dig.Container) error {
	err := container.Provide(func() (*Options, error) {
		return ProvideConfig()
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide gorm options")
	}

	err = container.Provide(func(opts *Options) (*gorm.DB, error) {
		return NewGorm(opts)
	})
	return errors.WrapIf(err, "failed to provide gorm db")
}
