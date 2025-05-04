package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	defaultLogger "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger/defaultlogger"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type Application struct {
	Container *dig.Container
	Echo      *echo.Echo
	Logger    logger.Logger
	Cfg       *config.Config
}

func NewApplication(container *dig.Container) *Application {
	app := &Application{}
	if err := container.Invoke(func(c *config.Config, e *echo.Echo, logger logger.Logger) error {
		app.Container = container
		app.Echo = e
		app.Logger = logger
		app.Cfg = c

		return nil
	}); err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}

	return app
}

func (a *Application) ResolveDependencyFunc(function interface{}) error {
	return a.Container.Invoke(function)
}

func (a *Application) ResolveRequiredDependencyFunc(function interface{}) {
	err := a.Container.Invoke(function)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve dependency: %v", err))
	}
}

func (a *Application) Run() {
	// https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
	// https://github.com/uber-go/fx/blob/master/app_test.go
	defaultDuration := time.Second * 20

	a.Start()
	<-a.Wait()

	stopCtx, stopCancellation := context.WithTimeout(context.Background(), defaultDuration)
	defer stopCancellation()
	a.Stop(stopCtx)
}

func (a *Application) Start() {
	echoStartHook(a)
}

func (a *Application) Stop(ctx context.Context) {
	echoStopHook(ctx, a)
	log.Println("Graceful shutdown complete.")
}

func (a *Application) Wait() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}

func echoStartHook(application *Application) {
	go func() {
		if err := application.Echo.Start(application.Cfg.EchoHTTPOptions.Port); !errors.Is(err, http.ErrServerClosed) {
			application.Logger.Fatalf("HTTP server error: %v", err)
		}
		application.Logger.Info("Stopped serving new HTTP connections.")
	}()
}

func echoStopHook(ctx context.Context, application *Application) {
	if err := application.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}
