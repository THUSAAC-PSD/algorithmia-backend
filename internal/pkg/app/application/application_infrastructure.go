package application

import (
	"strings"

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
	err := a.mapEndpoints()
	if err != nil {
		return errors.WrapIf(err, "failed to map endpoints")
	}

	err = a.migrateDatabase()
	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	err = a.seedProblemDifficulties()
	if err != nil {
		return errors.WrapIf(err, "failed to seed problem difficulties")
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
			&database.UserRole{},
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

func (a *Application) mustGenerateUUID(l logger.Logger) uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		l.Fatal("failed to generate UUID: ", err)
	}

	return id
}
