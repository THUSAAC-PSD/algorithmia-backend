package application

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/deletecontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/assigntesters"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/listmessage"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/listproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/markcomplete"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/reviewproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/sendmessage"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/testproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdifficulty/feature/listproblemdifficulty"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/listproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/submitproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/upsertproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

func (a *Application) AddHandlers() error {
	if err := a.Container.Provide(validator.New); err != nil {
		return errors.WrapIf(err, "failed to provide validator")
	}

	if err := a.Container.Provide(register.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide register command handler")
	}

	if err := a.Container.Provide(requestemailverification.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide request email verification command handler")
	}

	if err := a.Container.Provide(login.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide login command handler")
	}

	if err := a.Container.Provide(logout.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide logout command handler")
	}

	if err := a.Container.Provide(getcurrentuser.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide get current user query handler")
	}

	if err := a.Container.Provide(createcontest.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide create contest command handler")
	}

	if err := a.Container.Provide(deletecontest.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide delete contest command handler")
	}

	if err := a.Container.Provide(listcontest.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list contest query handler")
	}

	if err := a.Container.Provide(listproblemdifficulty.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list problem difficulty query handler")
	}

	if err := a.Container.Provide(upsertproblemdraft.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide upsert problem draft command handler")
	}

	if err := a.Container.Provide(listproblemdraft.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list problem draft query handler")
	}

	if err := a.Container.Provide(submitproblemdraft.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide submit problem draft command handler")
	}

	if err := a.Container.Provide(assigntesters.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide assign tester command handler")
	}

	if err := a.Container.Provide(listmessage.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list message query handler")
	}

	if err := a.Container.Provide(listproblem.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list problem query handler")
	}

	if err := a.Container.Provide(markcomplete.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide mark complete command handler")
	}

	if err := a.Container.Provide(reviewproblem.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide review problem command handler")
	}

	if err := a.Container.Provide(sendmessage.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide send message command handler")
	}

	if err := a.Container.Provide(testproblem.NewCommandHandler); err != nil {
		return errors.WrapIf(err, "failed to provide test problem command handler")
	}

	return nil
}
