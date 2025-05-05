package application

import (
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
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
		err := g.AutoMigrate(&database.User{}, &database.EmailVerificationCode{}, &database.Contest{})
		if err != nil {
			return err
		}

		return nil
	})
}
