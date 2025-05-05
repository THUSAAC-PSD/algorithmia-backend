package application

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
)

func (a *Application) ConfigMediator() error {
	return a.ResolveDependencyFunc(func(
		registerRepo register.Repository,
		requestEmailVerificationRepo requestemailverification.Repository,
		loginRepo login.Repository,
		createContestRepo createcontest.Repository,
		emailSender requestemailverification.EmailSender,
		passwordHasher register.PasswordHasher,
		passwordChecker login.PasswordChecker,
		loginSessionManager login.SessionManager,
		logoutSessionManager logout.SessionManager,
		uowFactory contract.UnitOfWorkFactory,
		l logger.Logger,
		authProvider contract.AuthProvider,
	) error {
		registerHandler := register.NewCommandHandler(registerRepo, passwordHasher, validator.New(), uowFactory, l)
		if err := mediatr.RegisterRequestHandler[*register.Command, *register.Response](registerHandler); err != nil {
			return errors.WrapIf(err, "failed to register register command handler")
		}

		requestEmailVerificationHandler := requestemailverification.NewCommandHandler(
			requestEmailVerificationRepo,
			emailSender,
			validator.New(),
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*requestemailverification.Command, mediatr.Unit](requestEmailVerificationHandler); err != nil {
			return errors.WrapIf(err, "failed to register request email verification command handler")
		}

		loginHandler := login.NewCommandHandler(
			loginRepo,
			passwordChecker,
			loginSessionManager,
			validator.New(),
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*login.Command, mediatr.Unit](loginHandler); err != nil {
			return errors.WrapIf(err, "failed to register login command handler")
		}

		logoutHandler := logout.NewCommandHandler(logoutSessionManager)
		if err := mediatr.RegisterRequestHandler[*logout.Command, mediatr.Unit](logoutHandler); err != nil {
			return errors.WrapIf(err, "failed to register logout command handler")
		}

		getCurrentUserHandler := getcurrentuser.NewQueryHandler(authProvider)
		if err := mediatr.RegisterRequestHandler[*getcurrentuser.Query, *getcurrentuser.Response](getCurrentUserHandler); err != nil {
			return errors.WrapIf(err, "failed to register get current user query handler")
		}

		createContestHandler := createcontest.NewCommandHandler(createContestRepo, validator.New())
		if err := mediatr.RegisterRequestHandler[*createcontest.Command, *createcontest.Response](createContestHandler); err != nil {
			return errors.WrapIf(err, "failed to register create contest command handler")
		}

		return nil
	})
}
