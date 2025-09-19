package echoweb

import (
	"net/http"
	"strings"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/log"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wader/gormstore/v2"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func AddEcho(container *dig.Container) error {
	if err := container.Provide(NewSessionAuthProvider,
		dig.As(new(contract.AuthProvider))); err != nil {
		return errors.WrapIf(err, "failed to provide session auth provider")
	}

	if err := container.Provide(func(l logger.Logger, opts *Options, db *gorm.DB) *echo.Echo {
		e := echo.New()

		e.HideBanner = true
		e.HTTPErrorHandler = httperror.Handler
		e.Validator = &customValidator{validator: validator.New()}

		e.Use(context.Middleware())
		e.Use(middleware.Recover())
		e.Use(log.EchoLogger(l))
		e.Use(middleware.BodyLimit(constant.BodyLimit))
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:5173", "https://algorithmia.thusaac.com", "http://algorithmia.thusaac.com"},
			AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
			AllowCredentials: true,
		}))
		e.Use(middleware.RequestID())
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: constant.GzipLevel,
			Skipper: func(c echo.Context) bool {
				if strings.Contains(c.Request().URL.Path, "ws/chat") &&
					c.Request().Header.Get("Upgrade") == "websocket" {
					return true
				}
				return false
			},
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
