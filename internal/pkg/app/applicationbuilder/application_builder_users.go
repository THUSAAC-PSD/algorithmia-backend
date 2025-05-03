package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (b *ApplicationBuilder) AddUsers() error {
	err := b.addRoutes()
	if err != nil {
		return err
	}

	err = b.addRepositories()
	if err != nil {
		return err
	}

	err = b.addPasswordHasher()
	if err != nil {
		return err
	}

	return nil
}

func (b *ApplicationBuilder) addRoutes() error {
	err := b.Container.Provide(func(e *echo.Echo, l logger.Logger) (*shared.UserEndpointParams, error) {
		v1 := e.Group("/api/v1")

		users := v1.Group("/users")
		auth := v1.Group("/auth")

		ep := &shared.UserEndpointParams{
			Logger:     l,
			Validator:  validator.New(),
			UsersGroup: users,
			AuthGroup:  auth,
		}

		return ep, nil
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide user endpoint params")
	}

	err = b.Container.Provide(func(ep *shared.UserEndpointParams) ([]contract.Endpoint, error) {
		registerEndpoint := register.NewEndpoint(ep)

		endpoints := []contract.Endpoint{registerEndpoint}
		return endpoints, nil
	})

	return errors.WrapIf(err, "failed to provide user endpoints")
}

func (b *ApplicationBuilder) addRepositories() error {
	err := b.Container.Provide(func(g *gorm.DB) (register.Repository, error) {
		return register.NewGormRepository(g), nil
	})

	return errors.WrapIf(err, "failed to provide user repositories")
}

func (b *ApplicationBuilder) addPasswordHasher() error {
	err := b.Container.Provide(func() (register.PasswordHasher, error) {
		return register.ArgonPasswordHasher{}, nil
	})

	return errors.WrapIf(err, "failed to provide password hasher")
}
