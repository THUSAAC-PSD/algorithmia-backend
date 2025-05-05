package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared/infrastructure"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func (b *ApplicationBuilder) AddUsers() error {
	if err := b.addRoutes(); err != nil {
		return err
	}

	if err := b.addRepositories(); err != nil {
		return err
	}

	if err := b.Container.Provide(infrastructure.NewArgonPasswordHasher,
		dig.As(new(register.PasswordHasher)),
		dig.As(new(login.PasswordChecker))); err != nil {
		return errors.WrapIf(err, "failed to provide argon password hasher")
	}

	if err := b.Container.Provide(requestemailverification.NewGomailEmailSender,
		dig.As(new(requestemailverification.EmailSender))); err != nil {
		return errors.WrapIf(err, "failed to provide gomail email sender")
	}

	if err := b.Container.Provide(infrastructure.NewHTTPSessionManager,
		dig.As(new(login.SessionManager)),
		dig.As(new(logout.SessionManager))); err != nil {
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
		logoutEndpoint := logout.NewEndpoint(ep)
		getCurrentUserEndpoint := getcurrentuser.NewEndpoint(ep)

		endpoints := []contract.Endpoint{
			registerEndpoint,
			requestEmailVerificationEndpoint,
			loginEndpoint,
			logoutEndpoint,
			getCurrentUserEndpoint,
		}
		return endpoints, nil
	})

	return errors.WrapIf(err, "failed to provide user endpoints")
}

func (b *ApplicationBuilder) addRepositories() error {
	if err := b.Container.Provide(register.NewGormRepository,
		dig.As(new(register.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide register repository")
	}

	if err := b.Container.Provide(requestemailverification.NewGormRepository,
		dig.As(new(requestemailverification.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide request email verification repository")
	}

	if err := b.Container.Provide(login.NewGormRepository,
		dig.As(new(login.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide login repository")
	}

	return nil
}
