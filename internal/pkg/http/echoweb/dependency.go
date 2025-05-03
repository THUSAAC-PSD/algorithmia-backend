package echoweb

import (
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/log"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

func AddEcho(container *dig.Container) error {
	err := container.Provide(func(l logger.Logger) *echo.Echo {
		e := echo.New()
		e.HideBanner = true

		skipper := func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "metrics") ||
				strings.Contains(c.Request().URL.Path, "health")
		}

		e.HTTPErrorHandler = httperror.Handler

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

		return e
	})
	return errors.WrapIf(err, "failed to provide echo instance")
}
