package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/assignproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/deletecontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listassignedproblems"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/unassignproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/assigntesters"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/checkoutdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/getproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/listmessage"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/listproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/markcomplete"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/reviewproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/sendmessage"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/testproblem"
	problemInfra "github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/infrastructure"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdifficulty"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdifficulty/feature/listproblemdifficulty"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/deleteproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/listproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/submitproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/upsertproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/listtester"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/manageuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/resetpassword"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/verifyemail"
	userInfra "github.com/THUSAAC-PSD/algorithmia-backend/internal/user/infrastructure"

	"emperror.dev/errors"
	"go.uber.org/dig"
)

func (b *ApplicationBuilder) AddFeatures() error {
	if err := b.addRoutes(); err != nil {
		return err
	}

	if err := b.addRepositories(); err != nil {
		return err
	}

	if err := b.Container.Provide(userInfra.NewArgonPasswordHasher,
		dig.As(new(register.PasswordHasher)),
		dig.As(new(requestemailverification.PasswordHasher)),
		dig.As(new(login.PasswordChecker)),
		dig.As(new(resetpassword.PasswordHasher))); err != nil {
		return errors.WrapIf(err, "failed to provide argon password hasher")
	}

	if err := b.Container.Provide(requestemailverification.NewPostmarkEmailSender,
		dig.As(new(requestemailverification.EmailSender))); err != nil {
		return errors.WrapIf(err, "failed to provide postmark email sender")
	}

	if err := b.Container.Provide(userInfra.NewHTTPSessionManager,
		dig.As(new(login.SessionManager)),
		dig.As(new(logout.SessionManager))); err != nil {
		return errors.WrapIf(err, "failed to provide http session manager")
	}

	return nil
}

func (b *ApplicationBuilder) addRoutes() error {
	// ======= Endpoint Params ========
	if err := b.Container.Provide(websocket.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide websocket endpoint params")
	}

	if err := b.Container.Provide(user.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide user endpoint params")
	}

	if err := b.Container.Provide(contest.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide contest endpoint params")
	}

	if err := b.Container.Provide(problem.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide problem endpoint params")
	}

	if err := b.Container.Provide(problemdraft.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide problem draft endpoint params")
	}

	if err := b.Container.Provide(problemdifficulty.NewEndpointParams); err != nil {
		return errors.WrapIf(err, "failed to provide problem difficulty endpoint params")
	}

	// ======== Endpoints ========
	if err := b.Container.Provide(websocket.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide websocket endpoint")
	}

	if err := b.Container.Provide(register.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide register endpoint")
	}

	if err := b.Container.Provide(requestemailverification.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide request email verification endpoint")
	}

	if err := b.Container.Provide(verifyemail.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide verify email endpoint")
	}

	if err := b.Container.Provide(login.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide login endpoint")
	}

	if err := b.Container.Provide(logout.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide logout endpoint")
	}

	if err := b.Container.Provide(getcurrentuser.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide get current user endpoint")
	}

	if err := b.Container.Provide(resetpassword.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide reset password endpoint")
	}

	if err := b.Container.Provide(manageuser.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide manage user endpoint")
	}

	if err := b.Container.Provide(createcontest.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide create contest endpoint")
	}

	if err := b.Container.Provide(listcontest.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list contest endpoint")
	}

	if err := b.Container.Provide(deletecontest.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide delete contest endpoint")
	}

	if err := b.Container.Provide(listproblemdifficulty.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list problem difficulty endpoint")
	}

	if err := b.Container.Provide(upsertproblemdraft.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide upsert problem draft endpoint")
	}

	if err := b.Container.Provide(listproblemdraft.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list problem draft endpoint")
	}

	if err := b.Container.Provide(submitproblemdraft.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide submit problem draft endpoint")
	}

	if err := b.Container.Provide(deleteproblemdraft.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide delete problem draft endpoint")
	}

	if err := b.Container.Provide(reviewproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide review problem endpoint")
	}

	if err := b.Container.Provide(testproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide test problem endpoint")
	}

	if err := b.Container.Provide(assigntesters.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide assign tester endpoint")
	}

	if err := b.Container.Provide(markcomplete.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide mark complete endpoint")
	}

	if err := b.Container.Provide(checkoutdraft.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide checkout draft endpoint")
	}

	if err := b.Container.Provide(listproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list problem endpoint")
	}

	if err := b.Container.Provide(getproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide get problem endpoint")
	}

	if err := b.Container.Provide(sendmessage.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide send message endpoint")
	}

	if err := b.Container.Provide(listmessage.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list message endpoint")
	}

	if err := b.Container.Provide(listtester.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list tester endpoint")
	}

	if err := b.Container.Provide(assignproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide assign problem endpoint")
	}

	if err := b.Container.Provide(unassignproblem.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide unassign problem endpoint")
	}

	if err := b.Container.Provide(listassignedproblems.NewEndpoint); err != nil {
		return errors.WrapIf(err, "failed to provide list assigned problems endpoint")
	}

	if err := b.Container.Provide(listassignedproblems.NewQueryHandler); err != nil {
		return errors.WrapIf(err, "failed to provide list assigned problems query handler")
	}

	if err := b.Container.Provide(func(
		websocketEndpoint *websocket.Endpoint,
		registerEndpoint *register.Endpoint,
		requestEmailVerificationEndpoint *requestemailverification.Endpoint,
		verifyEmailEndpoint *verifyemail.Endpoint,
		loginEndpoint *login.Endpoint,
		logoutEndpoint *logout.Endpoint,
		getCurrentUserEndpoint *getcurrentuser.Endpoint,
		resetPasswordEndpoint *resetpassword.Endpoint,
		createContestEndpoint *createcontest.Endpoint,
		listContestEndpoint *listcontest.Endpoint,
		deleteContestEndpoint *deletecontest.Endpoint,
		listAssignedProblemsEndpoint *listassignedproblems.Endpoint,
		listProblemDifficultyEndpoint *listproblemdifficulty.Endpoint,
		upsertProblemDraftEndpoint *upsertproblemdraft.Endpoint,
		listProblemDraftEndpoint *listproblemdraft.Endpoint,
		submitProblemDraftEndpoint *submitproblemdraft.Endpoint,
		deleteProblemDraftEndpoint *deleteproblemdraft.Endpoint,
		reviewProblemEndpoint *reviewproblem.Endpoint,
		testProblemEndpoint *testproblem.Endpoint,
		assignTesterEndpoint *assigntesters.Endpoint,
		markCompleteEndpoint *markcomplete.Endpoint,
		checkoutDraftEndpoint *checkoutdraft.Endpoint,
		listProblemEndpoint *listproblem.Endpoint,
		getProblemEndpoint *getproblem.Endpoint,
		sendMessageEndpoint *sendmessage.Endpoint,
		listMessageEndpoint *listmessage.Endpoint,
		listTesterEndpoint *listtester.Endpoint,
		assignProblemEndpoint *assignproblem.Endpoint,
		unassignProblemEndpoint *unassignproblem.Endpoint,
		manageUserEndpoint *manageuser.Endpoint,
	) []contract.Endpoint {
		return []contract.Endpoint{
			websocketEndpoint,
			registerEndpoint,
			requestEmailVerificationEndpoint,
			verifyEmailEndpoint,
			loginEndpoint,
			logoutEndpoint,
			getCurrentUserEndpoint,
			resetPasswordEndpoint,
			createContestEndpoint,
			listContestEndpoint,
			deleteContestEndpoint,
			listAssignedProblemsEndpoint,
			listProblemDifficultyEndpoint,
			upsertProblemDraftEndpoint,
			listProblemDraftEndpoint,
			submitProblemDraftEndpoint,
			deleteProblemDraftEndpoint,
			reviewProblemEndpoint,
			testProblemEndpoint,
			assignTesterEndpoint,
			markCompleteEndpoint,
			checkoutDraftEndpoint,
			listProblemEndpoint,
			getProblemEndpoint,
			sendMessageEndpoint,
			listMessageEndpoint,
			listTesterEndpoint,
			assignProblemEndpoint,
			unassignProblemEndpoint,
			manageUserEndpoint,
		}
	}); err != nil {
		return errors.WrapIf(err, "failed to provide endpoint array")
	}

	return nil
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

	if err := b.Container.Provide(verifyemail.NewGormRepository,
		dig.As(new(verifyemail.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide verify email repository")
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

	if err := b.Container.Provide(deletecontest.NewGormRepository,
		dig.As(new(deletecontest.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide delete contest repository")
	}

	if err := b.Container.Provide(listassignedproblems.NewGormRepository,
		dig.As(new(listassignedproblems.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list assigned problems repository")
	}

	if err := b.Container.Provide(listproblemdifficulty.NewGormRepository,
		dig.As(new(listproblemdifficulty.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list problem difficulty repository")
	}

	if err := b.Container.Provide(upsertproblemdraft.NewGormRepository,
		dig.As(new(upsertproblemdraft.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide upsert problem draft repository")
	}

	if err := b.Container.Provide(listproblemdraft.NewGormRepository,
		dig.As(new(listproblemdraft.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list problem draft repository")
	}

	if err := b.Container.Provide(submitproblemdraft.NewGormRepository,
		dig.As(new(submitproblemdraft.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide submit problem draft repository")
	}

	if err := b.Container.Provide(checkoutdraft.NewGormRepository,
		dig.As(new(checkoutdraft.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide checkout draft repository")
	}

	if err := b.Container.Provide(problemInfra.NewProblemActionGormRepository,
		dig.As(new(problemInfra.ProblemActionRepository)),
		dig.As(new(reviewproblem.Repository)),
		dig.As(new(testproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide shared problem action repository")
	}

	if err := b.Container.Provide(assigntesters.NewGormRepository,
		dig.As(new(assigntesters.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide assign tester repository")
	}

	if err := b.Container.Provide(markcomplete.NewGormRepository,
		dig.As(new(markcomplete.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide assign tester repository")
	}

	if err := b.Container.Provide(listproblem.NewGormRepository,
		dig.As(new(listproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list problem repository")
	}

	if err := b.Container.Provide(sendmessage.NewGormRepository,
		dig.As(new(sendmessage.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide send message repository")
	}

	if err := b.Container.Provide(listmessage.NewGormRepository,
		dig.As(new(listmessage.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list message repository")
	}

	if err := b.Container.Provide(getproblem.NewGormRepository,
		dig.As(new(getproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide get problem repository")
	}

	if err := b.Container.Provide(deleteproblemdraft.NewGormRepository,
		dig.As(new(deleteproblemdraft.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide delete problem draft repository")
	}

	if err := b.Container.Provide(listtester.NewGormRepository,
		dig.As(new(listtester.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide list tester repository")
	}

	if err := b.Container.Provide(assignproblem.NewGormRepository,
		dig.As(new(assignproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide assign problem repository")
	}

	if err := b.Container.Provide(unassignproblem.NewGormRepository,
		dig.As(new(unassignproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide unassign problem repository")
	}

	if err := b.Container.Provide(resetpassword.NewGormRepository,
		dig.As(new(resetpassword.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide reset password repository")
	}

	if err := b.Container.Provide(manageuser.NewGormRepository,
		dig.As(new(manageuser.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide manage user repository")
	}

	return nil
}
