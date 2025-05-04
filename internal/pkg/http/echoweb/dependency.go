package echoweb

import (
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/log"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

func AddEcho(container *dig.Container) error {
	if err := container.Provide(func() (*Options, error) {
		return ProvideConfig()
	}); err != nil {
		return errors.WrapIf(err, "failed to provide echo options")
	}

	err := container.Provide(func(l logger.Logger, opts *Options) *echo.Echo {
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
		e.Use(session.Middleware(sessions.NewCookieStore([]byte(opts.SessionSecret))))

		return e
	})
	return errors.WrapIf(err, "failed to provide echo instance")
}
