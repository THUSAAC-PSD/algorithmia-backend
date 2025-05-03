package application

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"

	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
)

func (a *Application) ConfigMediator() error {
	return a.ResolveDependencyFunc(func(
		registerRepo register.Repository,
		passwordHasher register.PasswordHasher,
	) error {
		registerHandler := register.NewCommandHandler(registerRepo, passwordHasher, validator.New())
		return mediatr.RegisterRequestHandler[*register.Command, *register.Response](registerHandler)
	})
}
