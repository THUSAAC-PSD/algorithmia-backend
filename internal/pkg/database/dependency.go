package database

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func AddGorm(container *dig.Container) error {
	if err := container.Provide(NewGorm); err != nil {
		return errors.WrapIf(err, "failed to provide gorm db")
	}

	if err := container.Provide(NewGormUnitOfWorkFactory, dig.As(new(contract.UnitOfWorkFactory))); err != nil {
		return errors.WrapIf(err, "failed to provide gorm unit of work factory")
	}

	return nil
}
