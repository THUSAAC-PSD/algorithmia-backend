package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/createcontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/deletecontest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/feature/listcontest"
	contestShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/contest/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/reviewproblem"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/feature/testproblem"
	problemShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared"
	problemInfra "github.com/THUSAAC-PSD/algorithmia-backend/internal/problem/shared/infrastructure"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdifficulty/feature/listproblemdifficulty"
	problemDifficultyShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdifficulty/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/listproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/submitproblemdraft"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/feature/upsertproblemdraft"
	problemDraftShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/getcurrentuser"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/logout"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/register"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/requestemailverification"
	userShared "github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"
	userInfra "github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared/infrastructure"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
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
		dig.As(new(login.PasswordChecker))); err != nil {
		return errors.WrapIf(err, "failed to provide argon password hasher")
	}

	if err := b.Container.Provide(requestemailverification.NewGomailEmailSender,
		dig.As(new(requestemailverification.EmailSender))); err != nil {
		return errors.WrapIf(err, "failed to provide gomail email sender")
	}

	if err := b.Container.Provide(userInfra.NewHTTPSessionManager,
		dig.As(new(login.SessionManager)),
		dig.As(new(logout.SessionManager))); err != nil {
		return errors.WrapIf(err, "failed to provide http session manager")
	}

	return nil
}

func (b *ApplicationBuilder) addRoutes() error {
	if err := b.Container.Provide(func(e *echo.Echo, v1Group *echoweb.V1Group) (*userShared.UserEndpointParams, error) {
		users := v1Group.Group.Group("/users")
		auth := v1Group.Group.Group("/auth")

		ep := &userShared.UserEndpointParams{
			Validator:  validator.New(),
			UsersGroup: users,
			AuthGroup:  auth,
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide user endpoint params")
	}

	if err := b.Container.Provide(func(e *echo.Echo, v1Group *echoweb.V1Group) (*contestShared.ContestEndpointParams, error) {
		contests := v1Group.Group.Group("/contests")

		ep := &contestShared.ContestEndpointParams{
			Validator:     validator.New(),
			ContestsGroup: contests,
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide contest endpoint params")
	}

	if err := b.Container.Provide(func(e *echo.Echo, v1Group *echoweb.V1Group) (*problemDifficultyShared.ProblemDifficultyEndpointParams, error) {
		problemDifficulties := v1Group.Group.Group("/problem-difficulties")

		ep := &problemDifficultyShared.ProblemDifficultyEndpointParams{
			ProblemDifficultiesGroup: problemDifficulties,
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide problem difficulty endpoint params")
	}

	if err := b.Container.Provide(func(e *echo.Echo, v1Group *echoweb.V1Group) (*problemDraftShared.ProblemDraftEndpointParams, error) {
		problemDrafts := v1Group.Group.Group("/problem-drafts")

		ep := &problemDraftShared.ProblemDraftEndpointParams{
			ProblemDraftsGroup: problemDrafts,
			Validator:          validator.New(),
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide problem draft endpoint params")
	}

	if err := b.Container.Provide(func(e *echo.Echo, v1Group *echoweb.V1Group) (*problemShared.ProblemEndpointParams, error) {
		problems := v1Group.Group.Group("/problems")

		ep := &problemShared.ProblemEndpointParams{
			ProblemsGroup: problems,
			Validator:     validator.New(),
		}

		return ep, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide problem endpoint params")
	}

	if err := b.Container.Provide(func(
		uep *userShared.UserEndpointParams,
		cep *contestShared.ContestEndpointParams,
		pdep *problemDifficultyShared.ProblemDifficultyEndpointParams,
		pdrep *problemDraftShared.ProblemDraftEndpointParams,
		pep *problemShared.ProblemEndpointParams,
	) ([]contract.Endpoint, error) {
		registerEndpoint := register.NewEndpoint(uep)
		requestEmailVerificationEndpoint := requestemailverification.NewEndpoint(uep)
		loginEndpoint := login.NewEndpoint(uep)
		logoutEndpoint := logout.NewEndpoint(uep)
		getCurrentUserEndpoint := getcurrentuser.NewEndpoint(uep)

		createContestEndpoint := createcontest.NewEndpoint(cep)
		listContestEndpoint := listcontest.NewEndpoint(cep)
		deleteContestEndpoint := deletecontest.NewEndpoint(cep)

		listProblemDifficultyEndpoint := listproblemdifficulty.NewEndpoint(pdep)

		upsertProblemDraftEndpoint := upsertproblemdraft.NewEndpoint(pdrep)
		listProblemDraftEndpoint := listproblemdraft.NewEndpoint(pdrep)
		submitProblemDraftEndpoint := submitproblemdraft.NewEndpoint(pdrep)

		reviewProblemEndpoint := reviewproblem.NewEndpoint(pep)
		testProblemEndpoint := testproblem.NewEndpoint(pep)

		endpoints := []contract.Endpoint{
			registerEndpoint,
			requestEmailVerificationEndpoint,
			loginEndpoint,
			logoutEndpoint,
			getCurrentUserEndpoint,

			createContestEndpoint,
			listContestEndpoint,
			deleteContestEndpoint,

			listProblemDifficultyEndpoint,

			upsertProblemDraftEndpoint,
			listProblemDraftEndpoint,
			submitProblemDraftEndpoint,

			reviewProblemEndpoint,
			testProblemEndpoint,
		}
		return endpoints, nil
	}); err != nil {
		return errors.WrapIf(err, "failed to provide endpoints")
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

	if err := b.Container.Provide(problemInfra.NewProblemActionGormRepository,
		dig.As(new(reviewproblem.Repository)),
		dig.As(new(testproblem.Repository))); err != nil {
		return errors.WrapIf(err, "failed to provide shared problem action repository")
	}

	return nil
}
