package application

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/deletecontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/assigntester"
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
	"github.com/mehdihadeli/go-mediatr"
)

func (a *Application) ConfigMediator() error {
	return a.ResolveDependencyFunc(func(
		registerRepo register.Repository,
		requestEmailVerificationRepo requestemailverification.Repository,
		loginRepo login.Repository,
		createContestRepo createcontest.Repository,
		listContestRepo listcontest.Repository,
		deleteContestRepo deletecontest.Repository,
		listProblemDifficultyRepo listproblemdifficulty.Repository,
		upsertProblemDraftRepo upsertproblemdraft.Repository,
		listProblemDraftRepo listproblemdraft.Repository,
		submitProblemDraftRepo submitproblemdraft.Repository,
		reviewProblemRepo reviewproblem.Repository,
		testProblemRepo testproblem.Repository,
		assignTesterRepo assigntester.Repository,
		markCompleteRepo markcomplete.Repository,
		listProblemRepo listproblem.Repository,
		sendMessageRepo sendmessage.Repository,
		emailSender requestemailverification.EmailSender,
		passwordHasher register.PasswordHasher,
		passwordChecker login.PasswordChecker,
		loginSessionManager login.SessionManager,
		logoutSessionManager logout.SessionManager,
		uowFactory contract.UnitOfWorkFactory,
		l logger.Logger,
		authProvider contract.AuthProvider,
		messageBroadcaster contract.MessageBroadcaster,
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

		listContestHandler := listcontest.NewQueryHandler(listContestRepo)
		if err := mediatr.RegisterRequestHandler[*listcontest.Query, *listcontest.Response](listContestHandler); err != nil {
			return errors.WrapIf(err, "failed to register list contest query handler")
		}

		deleteContestHandler := deletecontest.NewCommandHandler(deleteContestRepo, validator.New())
		if err := mediatr.RegisterRequestHandler[*deletecontest.Command, mediatr.Unit](deleteContestHandler); err != nil {
			return errors.WrapIf(err, "failed to register delete contest query handler")
		}

		listProblemDifficultyHandler := listproblemdifficulty.NewQueryHandler(listProblemDifficultyRepo)
		if err := mediatr.RegisterRequestHandler[*listproblemdifficulty.Query, *listproblemdifficulty.Response](listProblemDifficultyHandler); err != nil {
			return errors.WrapIf(err, "failed to register list problem difficulty query handler")
		}

		upsertProblemDraftHandler := upsertproblemdraft.NewCommandHandler(
			upsertProblemDraftRepo,
			validator.New(),
			authProvider,
		)
		if err := mediatr.RegisterRequestHandler[*upsertproblemdraft.Command, *upsertproblemdraft.Response](upsertProblemDraftHandler); err != nil {
			return errors.WrapIf(err, "failed to register upsert problem draft query handler")
		}

		listProblemDraftHandler := listproblemdraft.NewQueryHandler(listProblemDraftRepo, authProvider)
		if err := mediatr.RegisterRequestHandler[*listproblemdraft.Query, *listproblemdraft.Response](listProblemDraftHandler); err != nil {
			return errors.WrapIf(err, "failed to register list problem draft query handler")
		}

		submitProblemDraftHandler := submitproblemdraft.NewCommandHandler(
			submitProblemDraftRepo,
			validator.New(),
			authProvider,
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*submitproblemdraft.Command, *submitproblemdraft.Response](submitProblemDraftHandler); err != nil {
			return errors.WrapIf(err, "failed to register submit problem draft command handler")
		}

		reviewProblemHandler := reviewproblem.NewCommandHandler(
			reviewProblemRepo,
			validator.New(),
			authProvider,
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*reviewproblem.Command, *reviewproblem.Response](reviewProblemHandler); err != nil {
			return errors.WrapIf(err, "failed to register review problem command handler")
		}

		testProblemHandler := testproblem.NewCommandHandler(
			testProblemRepo,
			validator.New(),
			authProvider,
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*testproblem.Command, *testproblem.Response](testProblemHandler); err != nil {
			return errors.WrapIf(err, "failed to register test problem command handler")
		}

		assignTesterHandler := assigntester.NewCommandHandler(
			assignTesterRepo,
			validator.New(),
			authProvider,
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*assigntester.Command, mediatr.Unit](assignTesterHandler); err != nil {
			return errors.WrapIf(err, "failed to register assign tester command handler")
		}

		markCompleteHandler := markcomplete.NewCommandHandler(
			markCompleteRepo,
			validator.New(),
			authProvider,
			uowFactory,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*markcomplete.Command, mediatr.Unit](markCompleteHandler); err != nil {
			return errors.WrapIf(err, "failed to register mark complete command handler")
		}

		listProblemHandler := listproblem.NewQueryHandler(
			listProblemRepo,
			authProvider,
		)
		if err := mediatr.RegisterRequestHandler[*listproblem.Query, *listproblem.Response](listProblemHandler); err != nil {
			return errors.WrapIf(err, "failed to register list problem query handler")
		}

		sendMessageHandler := sendmessage.NewCommandHandler(
			sendMessageRepo,
			validator.New(),
			authProvider,
			uowFactory,
			messageBroadcaster,
			l,
		)
		if err := mediatr.RegisterRequestHandler[*sendmessage.Command, mediatr.Unit](sendMessageHandler); err != nil {
			return errors.WrapIf(err, "failed to register send message command handler")
		}

		return nil
	})
}
