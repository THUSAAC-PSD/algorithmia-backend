package echoweb

import (
	"strings"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/log"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wader/gormstore/v2"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func AddEcho(container *dig.Container) error {
	if err := container.Provide(ProvideConfig); err != nil {
		return errors.WrapIf(err, "failed to provide echo options")
	}

	if err := container.Provide(NewSessionAuthProvider,
		dig.As(new(contract.AuthProvider))); err != nil {
		return errors.WrapIf(err, "failed to provide session auth provider")
	}

	if err := container.Provide(func(l logger.Logger, opts *Options, db *gorm.DB) *echo.Echo {
		e := echo.New()
		e.HideBanner = true

		skipper := func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "metrics") ||
				strings.Contains(c.Request().URL.Path, "health")
		}

		e.HTTPErrorHandler = httperror.Handler

		e.Use(context.Middleware())
		e.Use(middleware.Recover())
		e.Use(
			log.EchoLogger(
				l,
				log.WithSkipper(skipper),
			),
		)
		e.Use(middleware.BodyLimit(constant.BodyLimit))
		e.Use(middleware.RequestID())
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level:   constant.GzipLevel,
			Skipper: skipper,
		}))

		store := gormstore.New(db, []byte(opts.SessionSecret))

		quit := make(chan struct{})
		go store.PeriodicCleanup(1*time.Hour, quit)

		e.Use(session.Middleware(store))

		return e
	}); err != nil {
		return errors.WrapIf(err, "failed to provide echo instance")
	}

	if err := container.Provide(NewV1Group); err != nil {
		return errors.WrapIf(err, "failed to provide v1 group")
	}

	return nil
}
