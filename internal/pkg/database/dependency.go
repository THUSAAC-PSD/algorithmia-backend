package database

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database/options"

	"emperror.dev/errors"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func AddGorm(container *dig.Container) error {
	err := container.Provide(func() (*options.GormOptions, error) {
		return options.ProvideConfig()
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide gorm options")
	}

	err = container.Provide(func(opts *options.GormOptions) (*gorm.DB, error) {
		return NewGorm(opts)
	})
	return errors.WrapIf(err, "failed to provide gorm db")
}
