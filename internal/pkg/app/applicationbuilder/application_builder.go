package applicationbuilder

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/app/application"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type ApplicationBuilder struct {
	Container   *dig.Container
	Logger      logger.Logger
	Environment environment.Environment
	overrides   []*Override
}

type Override struct {
	DecoratorFunc interface{}
	Opts          []dig.DecorateOption
}

func NewApplicationBuilder(environments ...environment.Environment) *ApplicationBuilder {
	container := dig.New()

	err := logger.AddLogger(container)
	if err != nil {
		log.Fatalln(err)
	}

	setConfigPath()
	err = environment.AddEnv(container, environments...)
	if err != nil {
		log.Fatalln(err)
	}

	var l logger.Logger
	var env environment.Environment

	err = container.Invoke(func(logger logger.Logger, environment environment.Environment) error {
		env = environment
		l = logger

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	appBuilder := &ApplicationBuilder{Container: container, Logger: l, Environment: env}

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

func setConfigPath() {
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	pn := viper.Get(constant.ProjectNameEnv)
	if pn == nil {
		return
	}
	for !strings.HasSuffix(wd, pn.(string)) {
		wd = filepath.Dir(wd)
	}

	absCurrentDir, _ := filepath.Abs(wd)
	viper.Set(constant.AppRootPath, absCurrentDir)

	configPath := filepath.Join(absCurrentDir, "config")
	viper.Set(constant.ConfigPath, configPath)
}
