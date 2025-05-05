package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listcontest"
	contestShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"
	userShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared/infrastructure"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func (b *ApplicationBuilder) AddUsers() error {
	if err := b.addUserRoutes(); err != nil {
		return err
	}

	if err := b.addUserRepositories(); err != nil {
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

func (b *ApplicationBuilder) addUserRoutes() error {
	if err := b.Container.Provide(func(e *echo.Echo, l logger.Logger, v1Group *echoweb.V1Group) (*userShared.UserEndpointParams, error) {
		users := v1Group.Group.Group("/users")
		auth := v1Group.Group.Group("/auth")

		ep := &userShared.UserEndpointParams{
			Logger:     l,
			Validator:  validator.New(),
			UsersGroup: users,
			AuthGroup:  auth,
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide user endpoint params")
	}

	if err := b.Container.Provide(func(e *echo.Echo, l logger.Logger, v1Group *echoweb.V1Group) (*contestShared.ContestEndpointParams, error) {
		contests := v1Group.Group.Group("/contests")

		ep := &contestShared.ContestEndpointParams{
			Logger:        l,
			Validator:     validator.New(),
			ContestsGroup: contests,
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide contest endpoint params")
	}

	if err := b.Container.Provide(func(uep *userShared.UserEndpointParams, cep *contestShared.ContestEndpointParams) ([]contract.Endpoint, error) {
		registerEndpoint := register.NewEndpoint(uep)
		requestEmailVerificationEndpoint := requestemailverification.NewEndpoint(uep)
		loginEndpoint := login.NewEndpoint(uep)
		logoutEndpoint := logout.NewEndpoint(uep)
		getCurrentUserEndpoint := getcurrentuser.NewEndpoint(uep)

		createContestEndpoint := createcontest.NewEndpoint(cep)
		listContestEndpoint := listcontest.NewEndpoint(cep)

		endpoints := []contract.Endpoint{
			registerEndpoint,
			requestEmailVerificationEndpoint,
			loginEndpoint,
			logoutEndpoint,
			getCurrentUserEndpoint,

			createContestEndpoint,
			listContestEndpoint,
		}
		return endpoints, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide endpoints")
	}

	return nil
}

func (b *ApplicationBuilder) addUserRepositories() error {
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

	if err := b.Container.Provide(createcontest.NewGormRepository,
		dig.As(new(createcontest.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide create contest repository")
	}

	if err := b.Container.Provide(listcontest.NewGormRepository,
		dig.As(new(listcontest.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list contest repository")
	}

	return nil
}
