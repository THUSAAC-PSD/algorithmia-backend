package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
	"github.com/caarlos0/env/v8"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

func BindConfig[T any](environments ...environment.Environment) (T, error) {
	return BindConfigKey[T]("", environments...)
}

func BindConfigKey[T any](configKey string, environments ...environment.Environment) (T, error) {
	var configPath string

	e := environment.Environment("")
	if len(environments) > 0 {
		e = environments[0]
	} else {
		e = constant.Dev
	}

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	configPathFromEnv := viper.Get(constant.ConfigPath)
	if configPathFromEnv != nil {
		configPath = configPathFromEnv.(string)
	} else {
		// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
		// https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
		d, err := getConfigRootPath()
		if err != nil {
			return *new(T), err
		}

		configPath = d
	}

	cfg := typemapper.GenericInstanceByT[T]()

	// https://github.com/spf13/viper/issues/390#issuecomment-718756752
	viper.SetConfigName(fmt.Sprintf("config.%s", e))
	viper.AddConfigPath(configPath)
	viper.SetConfigType(constant.JSON)

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), errors.WrapIf(err, "failed to read config file")
	}

	if len(configKey) == 0 {
		if err := viper.Unmarshal(cfg); err != nil {
			return *new(T), errors.WrapIf(err, "failed to unmarshal config")
		}
	} else {
		if err := viper.UnmarshalKey(configKey, cfg); err != nil {
			return *new(T), errors.WrapIf(err, "failed to unmarshal config")
		}
	}

	viper.AutomaticEnv()

	// https://github.com/caarlos0/env
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// https://github.com/mcuadros/go-defaults
	defaults.SetDefaults(cfg)

	return cfg, nil
}

func getConfigRootPath() (string, error) {
	// Getwd gives us the current working directory that we are running our app with `go run` command. in goland we can specify this working directory for the project
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fmt.Printf("Current working directory is: %s\n", wd)

	// Get the absolute path of the executed project directory
	absCurrentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(absCurrentDir, "config")
	return configPath, nil
}
