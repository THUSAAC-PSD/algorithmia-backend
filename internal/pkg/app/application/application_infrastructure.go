package application

import (
	"strings"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (a *Application) ConfigInfrastructure() error {
	if err := a.mapEndpoints(); err != nil {
		return errors.WrapIf(err, "failed to map endpoints")
	}

	if err := a.migrateDatabase(); err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	if err := a.seedProblemDifficulties(); err != nil {
		return errors.WrapIf(err, "failed to seed problem difficulties")
	}

	if err := a.seedPermissions(); err != nil {
		return errors.WrapIf(err, "failed to seed permissions")
	}

	if err := a.seedRoles(); err != nil {
		return errors.WrapIf(err, "failed to seed roles")
	}

	return nil
}

func (a *Application) mapEndpoints() error {
	a.ResolveRequiredDependencyFunc(func(endpoints []contract.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	})

	a.ResolveRequiredDependencyFunc(func(e *echo.Echo, l logger.Logger) {
		l.Info("Registered routes:")
		for _, route := range e.Routes() {
			name, _ := strings.CutPrefix(route.Name, "github.com/THUSAAC-PSD/algorithmia-backend/internal/")
			l.Infof("%s %s: %s", route.Method, route.Path, name)
		}
	})

	return nil
}

func (a *Application) migrateDatabase() error {
	return a.ResolveDependencyFunc(func(g *gorm.DB) error {
		err := g.AutoMigrate(
			&database.User{},
			&database.Role{},
			&database.Permission{},
			&database.EmailVerificationCode{},
			&database.Contest{},
			&database.ProblemDifficulty{},
			&database.ProblemDifficultyDisplayName{},
			&database.ProblemDraft{},
			&database.ProblemDraftDetail{},
			&database.ProblemDraftExample{},
			&database.Problem{},
			&database.ProblemVersion{},
			&database.ProblemVersionDetail{},
			&database.ProblemVersionExample{},
			&database.ProblemReview{},
			&database.ProblemTestResult{},
			&database.Media{},
			&database.ProblemChatMessage{},
			&database.ProblemChatMessageAttachment{},
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (a *Application) seedProblemDifficulties() error {
	return a.ResolveDependencyFunc(func(g *gorm.DB, l logger.Logger) error {
		var count int64
		if err := g.Model(&database.ProblemDifficulty{}).Count(&count).Error; err != nil {
			return errors.WrapIf(err, "failed to count problem difficulties")
		}

		if count > 0 {
			return nil
		}

		problemDifficulties := []database.ProblemDifficulty{
			{ProblemDifficultyID: a.mustGenerateUUID(l)},
			{ProblemDifficultyID: a.mustGenerateUUID(l)},
			{ProblemDifficultyID: a.mustGenerateUUID(l)},
		}

		problemDifficulties[0].DisplayNames = []database.ProblemDifficultyDisplayName{
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[0].ProblemDifficultyID,
				DisplayName:         "Easy",
				Language:            constant.LanguageEnUS,
			},
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[0].ProblemDifficultyID,
				DisplayName:         "简单",
				Language:            constant.LanguageZhCN,
			},
		}

		problemDifficulties[1].DisplayNames = []database.ProblemDifficultyDisplayName{
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[1].ProblemDifficultyID,
				DisplayName:         "Medium",
				Language:            constant.LanguageEnUS,
			},
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[1].ProblemDifficultyID,
				DisplayName:         "中等",
				Language:            constant.LanguageZhCN,
			},
		}

		problemDifficulties[2].DisplayNames = []database.ProblemDifficultyDisplayName{
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[2].ProblemDifficultyID,
				DisplayName:         "Hard",
				Language:            constant.LanguageEnUS,
			},
			{
				DisplayNameID:       a.mustGenerateUUID(l),
				ProblemDifficultyID: problemDifficulties[2].ProblemDifficultyID,
				DisplayName:         "困难",
				Language:            constant.LanguageZhCN,
			},
		}

		if err := g.Create(&problemDifficulties).Error; err != nil {
			return errors.WrapIf(err, "failed to create problem difficulties")
		}

		return nil
	})
}

func (a *Application) seedPermissions() error {
	return a.ResolveDependencyFunc(func(g *gorm.DB, l logger.Logger) error {
		var count int64
		if err := g.Model(&database.Permission{}).Count(&count).Error; err != nil {
			return errors.WrapIf(err, "failed to count permissions")
		}

		if count > 0 {
			return nil
		}

		now := time.Now()

		permissions := make([]database.Permission, 0)
		addPermission := func(name, description string) {
			permissions = append(permissions, database.Permission{
				PermissionID: a.mustGenerateUUID(l),
				Name:         name,
				Description:  description,
				CreatedAt:    now,
			})
		}

		addPermission(constant.PermissionMediaUploadForDraftOwn, "Upload media for one's own problem draft")
		addPermission(constant.PermissionMediaUploadForChatOwn, "Upload media for one's own problem chat")

		addPermission(constant.PermissionUserListAll, "List all users")
		addPermission(constant.PermissionUserReadProfileAny, "Read any user's profile")
		addPermission(constant.PermissionUserReadProfileOwn, "Read one's own profile")
		addPermission(constant.PermissionUserUpdateProfileOwn, "Update one's own profile")
		addPermission(constant.PermissionUserManageRolesAny, "Manage any user's roles (assign/remove)")

		addPermission(constant.PermissionRoleListAll, "List all roles")
		addPermission(constant.PermissionRoleManageAny, "Manage any role (create/update/delete)")

		addPermission(constant.PermissionProblemDraftCreate, "Create a problem draft")
		addPermission(constant.PermissionProblemDraftReadOwn, "Read one's own problem draft")
		addPermission(constant.PermissionProblemDraftUpdateOwn, "Update one's own problem draft")
		addPermission(constant.PermissionProblemDraftDeleteOwn, "Delete one's own problem draft")
		addPermission(constant.PermissionProblemDraftSubmitOwn, "Submit one's own problem draft")

		addPermission(constant.PermissionProblemListAll, "List all problems")
		addPermission(constant.PermissionProblemListCreatedOwn, "List problems created by oneself")
		addPermission(constant.PermissionProblemListAwaitingReviewAll, "List all problems awaiting review")
		addPermission(constant.PermissionProblemListAssignedTest, "List problems assigned to oneself as a tester")
		addPermission(constant.PermissionProblemReadDetailsAny, "Read details of any problem")
		addPermission(constant.PermissionProblemReadDetailsCreatedOwn, "Read details of problems created by oneself")
		addPermission(
			constant.PermissionProblemReadDetailsAwaitingReviewAny,
			"Read details of any problems awaiting review",
		)
		addPermission(
			constant.PermissionProblemReadDetailsAssignedTest,
			"Read details of problems assigned to oneself as a tester",
		)
		addPermission(constant.PermissionProblemReviewAny, "Review any problem")
		addPermission(constant.PermissionProblemReviewOverride, "Review any problem, even if not the reviewer")
		addPermission(constant.PermissionProblemAssignTesters, "Assign testers to problems")
		addPermission(
			constant.PermissionProblemTestAssigned,
			"Submit test result to problems assigned to oneself as a tester",
		)
		addPermission(
			constant.PermissionProblemTestOverride,
			"Submit test result to any problem, even if not assigned as a tester",
		)

		addPermission(constant.PermissionContestListAll, "List all contests")
		addPermission(constant.PermissionContestReadDetailsAny, "Read details of any contest")

		addPermission(constant.PermissionContestCreate, "Create a contest")
		addPermission(constant.PermissionContestUpdateAny, "Update any contest")
		addPermission(constant.PermissionContestDeleteAny, "Delete any contest")

		addPermission(constant.PermissionContestAssignProblemAny, "Assign problems to any contest")
		addPermission(constant.PermissionContestUnassignProblemAny, "Unassign problems from any contest")

		if err := g.Create(&permissions).Error; err != nil {
			return errors.WrapIf(err, "failed to create permissions")
		}

		return nil
	})
}

func (a *Application) seedRoles() error {
	return a.ResolveDependencyFunc(func(g *gorm.DB, l logger.Logger) error {
		var count int64
		if err := g.Model(&database.Role{}).Count(&count).Error; err != nil {
			return errors.WrapIf(err, "failed to count roles")
		}

		if count > 0 {
			return nil
		}

		now := time.Now()

		roles := make([]database.Role, 0)
		addRole := func(name, description string, permissionNames []string, isSuperAdmin bool) {
			permissions := make([]database.Permission, 0, len(permissionNames))
			for _, permissionName := range permissionNames {
				permissions = append(permissions, database.Permission{Name: permissionName})
			}

			roles = append(roles, database.Role{
				RoleID:       a.mustGenerateUUID(l),
				Name:         name,
				Description:  description,
				IsSuperAdmin: isSuperAdmin,
				Permissions:  &permissions,
				CreatedAt:    now,
			})
		}

		addRole("super_admin", "Super admin", []string{}, true)

		addRole("setter", "Problem setter", []string{
			constant.PermissionMediaUploadForChatOwn,
			constant.PermissionMediaUploadForDraftOwn,
			constant.PermissionUserReadProfileOwn,
			constant.PermissionUserUpdateProfileOwn,
			constant.PermissionProblemDraftCreate,
			constant.PermissionProblemDraftReadOwn,
			constant.PermissionProblemDraftUpdateOwn,
			constant.PermissionProblemDraftDeleteOwn,
			constant.PermissionProblemDraftSubmitOwn,
			constant.PermissionProblemListCreatedOwn,
			constant.PermissionProblemReadDetailsCreatedOwn,
			constant.PermissionContestListAll,
			constant.PermissionContestReadDetailsAny,
		}, false)

		addRole("reviewer", "Problem reviewer", []string{
			constant.PermissionMediaUploadForChatOwn,
			constant.PermissionUserReadProfileOwn,
			constant.PermissionUserUpdateProfileOwn,
			constant.PermissionProblemListAwaitingReviewAll,
			constant.PermissionProblemReadDetailsAwaitingReviewAny,
			constant.PermissionProblemReviewAny,
			constant.PermissionProblemAssignTesters,
		}, false)

		addRole("tester", "Problem tester", []string{
			constant.PermissionMediaUploadForChatOwn,
			constant.PermissionUserReadProfileOwn,
			constant.PermissionUserUpdateProfileOwn,
			constant.PermissionProblemListAssignedTest,
			constant.PermissionProblemReadDetailsAssignedTest,
			constant.PermissionProblemTestAssigned,
		}, false)

		addRole("contest_manager", "Contest manager", []string{
			constant.PermissionUserReadProfileOwn,
			constant.PermissionUserUpdateProfileOwn,
			constant.PermissionContestCreate,
			constant.PermissionContestListAll,
			constant.PermissionContestReadDetailsAny,
			constant.PermissionContestUpdateAny,
			constant.PermissionContestDeleteAny,
			constant.PermissionContestAssignProblemAny,
			constant.PermissionContestUnassignProblemAny,
			constant.PermissionProblemListAll,
		}, false)

		permissionNames := make([]string, 0)
		for _, role := range roles {
			for _, permission := range *role.Permissions {
				permissionNames = append(permissionNames, permission.Name)
			}
		}

		var permissions []database.Permission
		if err := g.Model(&database.Permission{}).
			Where("name IN ?", permissionNames).
			Scan(&permissions).Error; err != nil {
			return errors.WrapIf(err, "failed to get permissions")
		}

		for i, role := range roles {
			for j, permission := range *role.Permissions {
				for _, p := range permissions {
					if permission.Name == p.Name {
						(*roles[i].Permissions)[j] = p
						break
					}
				}
			}
		}

		if err := g.Create(&roles).Error; err != nil {
			return errors.WrapIf(err, "failed to create roles")
		}

		return nil
	})
}

func (a *Application) mustGenerateUUID(l logger.Logger) uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		l.Fatal("failed to generate UUID: ", err)
	}

	return id
}
