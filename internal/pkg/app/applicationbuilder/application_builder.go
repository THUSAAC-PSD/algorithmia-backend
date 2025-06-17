package applicationbuilder

import (
	"log"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/app/application"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"go.uber.org/dig"
)

type ApplicationBuilder struct {
	Container *dig.Container
	Logger    logger.Logger
	overrides []*Override
}

type Override struct {
	DecoratorFunc interface{}
	Opts          []dig.DecorateOption
}

func NewApplicationBuilder() *ApplicationBuilder {
	container := dig.New()

	if err := container.Provide(config.BindAllConfigs); err != nil {
		log.Fatalln(err)
	}

	if err := container.Provide(func(cfg *config.Config) *environment.Environment { return &cfg.Environment }); err != nil {
		log.Fatal(err)
	}

	if err := container.Provide(func(cfg *config.Config) *logger.Options { return &cfg.LoggerOptions }); err != nil {
		log.Fatal(err)
	}

	err := logger.AddLogger(container)
	if err != nil {
		log.Fatalln(err)
	}

	var l logger.Logger
	err = container.Invoke(func(logger logger.Logger) error {
		l = logger

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	appBuilder := &ApplicationBuilder{Container: container, Logger: l}

	return appBuilder
}

func (b *ApplicationBuilder) Build() *application.Application {
	for _, override := range b.overrides {
		err := b.Container.Decorate(override.DecoratorFunc, override.Opts...)
		if err != nil {
			b.Logger.Fatal(err)
		}
	}

	container := b.Container
	app := application.NewApplication(container)

	return app
}

func (b *ApplicationBuilder) WithOverride(decoratorFunc interface{}, opts ...dig.DecorateOption) *ApplicationBuilder {
	b.overrides = append(b.overrides, &Override{decoratorFunc, opts})
	return b
}
