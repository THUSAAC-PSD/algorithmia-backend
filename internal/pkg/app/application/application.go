package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	defaultLogger "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger/defaultlogger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type Application struct {
	Container    *dig.Container
	Echo         *echo.Echo
	WebsocketHub *websocket.Hub
	Logger       logger.Logger
	EchoOptions  *echoweb.Options

	appCtx    context.Context
	appCancel context.CancelFunc
}

func NewApplication(container *dig.Container) *Application {
	appCtx, appCancel := context.WithCancel(context.Background())

	app := &Application{
		appCtx:    appCtx,
		appCancel: appCancel,
	}

	if err := container.Invoke(func(opts *echoweb.Options, e *echo.Echo, wh *websocket.Hub, logger logger.Logger) error {
		app.Container = container
		app.Echo = e
		app.WebsocketHub = wh
		app.Logger = logger
		app.EchoOptions = opts

		return nil
	}); err != nil {
		l := defaultLogger.GetLogger()
		if app.Logger != nil {
			l = app.Logger
		}

		l.Fatal(err)
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

	a.Logger.Info("Starting application...")
	a.Start()

	<-a.Wait()
	a.Logger.Info("Received shutdown signal, shutting down...")

	stopCtx, stopCancellation := context.WithTimeout(context.Background(), defaultDuration)
	defer stopCancellation()

	a.Stop(stopCtx)
	a.Logger.Info("Application stopped")
}

func (a *Application) Start() {
	go func() {
		if err := a.Echo.Start(a.EchoOptions.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Fatalf("HTTP server error: %v", err)
		}

		a.Logger.Info("Stopped serving new HTTP connections.")
	}()

	go func() {
		a.Logger.Info("WebSocket Hub starting...")
		a.WebsocketHub.Run(a.appCtx)
		a.Logger.Info("WebSocket Hub stopped.")
	}()
}

func (a *Application) Stop(ctx context.Context) {
	a.Logger.Info("Stopping application components...")

	if a.appCancel != nil {
		a.appCancel()
	}

	time.Sleep(1 * time.Second) // Give hub a moment to start closing connections

	a.Logger.Info("Shutting down HTTP server...")
	if err := a.Echo.Shutdown(ctx); err != nil {
		a.Logger.Errorf("HTTP server shutdown error: %v", err)
	} else {
		a.Logger.Info("HTTP server shutdown complete.")
	}
}

func (a *Application) Wait() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}
