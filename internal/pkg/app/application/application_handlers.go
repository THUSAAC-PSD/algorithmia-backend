package application

import (
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
		emailSender requestemailverification.EmailSender,
		passwordHasher register.PasswordHasher,
	) error {
		registerHandler := register.NewCommandHandler(registerRepo, passwordHasher, validator.New())
		if err := mediatr.RegisterRequestHandler[*register.Command, *register.Response](registerHandler); err != nil {
			return errors.WrapIf(err, "failed to register register command handler")
		}

		requestEmailVerificationHandler := requestemailverification.NewCommandHandler(
			requestEmailVerificationRepo,
			emailSender,
			validator.New(),
		)
		if err := mediatr.RegisterRequestHandler[*requestemailverification.Command, mediatr.Unit](requestEmailVerificationHandler); err != nil {
			return errors.WrapIf(err, "failed to register request email verification command handler")
		}

		return nil
	})
}
