package logger

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

var optionName = strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[Options]())

type Options struct {
	LogLevel      string `mapstructure:"level"`
	CallerEnabled bool   `mapstructure:"callerEnabled"`
}

func ProvideConfig(env environment.Environment) (*Options, error) {
	return config.BindConfigKey[*Options](optionName, env)
}
