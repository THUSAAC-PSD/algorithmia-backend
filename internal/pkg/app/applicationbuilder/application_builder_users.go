package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/gomail"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared/infrastructure"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (b *ApplicationBuilder) AddUsers() error {
	if err := b.addRoutes(); err != nil {
		return err
	}

	if err := b.addRepositories(); err != nil {
		return err
	}

	if err := b.Container.Provide(func() (register.PasswordHasher, error) {
		return infrastructure.ArgonPasswordHasher{}, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide password hasher")
	}

	if err := b.Container.Provide(func() (login.PasswordChecker, error) {
		return infrastructure.ArgonPasswordHasher{}, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide password checker")
	}

	if err := b.Container.Provide(func(opts *gomail.Options) (requestemailverification.EmailSender, error) {
		return requestemailverification.NewGomailEmailSender(opts)
	}); err != nil {
		return errors.WrapIf(err, "failed to provide gomail email sender")
	}

	if err := b.Container.Provide(func() login.SessionManager {
		return login.NewHTTPSessionManager()
	}); err != nil {
		return errors.WrapIf(err, "failed to provide http session manager")
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
		requestEmailVerificationEndpoint := requestemailverification.NewEndpoint(ep)
		loginEndpoint := login.NewEndpoint(ep)

		endpoints := []contract.Endpoint{registerEndpoint, requestEmailVerificationEndpoint, loginEndpoint}
		return endpoints, nil
	})

	return errors.WrapIf(err, "failed to provide user endpoints")
}

func (b *ApplicationBuilder) addRepositories() error {
	err := b.Container.Provide(func(g *gorm.DB) (register.Repository, error) {
		return register.NewGormRepository(g), nil
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide register repository")
	}

	err = b.Container.Provide(func(g *gorm.DB) (requestemailverification.Repository, error) {
		return requestemailverification.NewGormRepository(g), nil
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide request email verification repository")
	}

	err = b.Container.Provide(func(g *gorm.DB) (login.Repository, error) {
		return login.NewGormRepository(g), nil
	})
	if err != nil {
		return errors.WrapIf(err, "failed to provide login repository")
	}

	return nil
}
