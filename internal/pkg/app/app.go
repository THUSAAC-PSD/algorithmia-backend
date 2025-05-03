package app

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/app/application"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/app/applicationbuilder"
)

type App struct{}

func NewApp() *App {
	app := &App{}
	return app
}

func (a *App) Run() {
	builder := createApplicationBuilder()

	app := builder.Build()

	configureApplication(app)

	app.Run()
}

func configureApplication(app *application.Application) {
	if err := app.ConfigMediator(); err != nil {
		app.Logger.Fatal(err)
	}

	if err := app.ConfigInfrastructure(); err != nil {
		app.Logger.Fatal(err)
	}
}

func createApplicationBuilder() *applicationbuilder.ApplicationBuilder {
	builder := applicationbuilder.NewApplicationBuilder()

	builder.AddInfrastructure()

	if err := builder.AddUsers(); err != nil {
		builder.Logger.Fatal(err)
	}

	return builder
}
