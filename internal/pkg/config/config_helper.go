package config

import (
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func BindAllConfigs() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		// Log the error but do not return it, as .env file is optional
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.BindEnv("environment", "APP_ENV")

	// LoggerOptions
	_ = viper.BindEnv("loggerOptions.level", "LOG_LEVEL")
	_ = viper.BindEnv("loggerOptions.callerEnabled", "LOG_CALLER_ENABLED")

	// GormOptions
	_ = viper.BindEnv("gormOptions.host", "DB_HOST")
	_ = viper.BindEnv("gormOptions.port", "DB_PORT")
	_ = viper.BindEnv("gormOptions.user", "DB_USER")
	_ = viper.BindEnv("gormOptions.password", "DB_PASSWORD")
	_ = viper.BindEnv("gormOptions.dbName", "DB_NAME")
	_ = viper.BindEnv("gormOptions.sslMode", "DB_SSL_MODE")
	_ = viper.BindEnv("gormOptions.useInMemory", "DB_USE_IN_MEMORY")

	// EchoHttpOptions
	_ = viper.BindEnv("echoHttpOptions.port", "PORT")
	_ = viper.BindEnv("echoHttpOptions.sessionSecret", "SESSION_SECRET")

	// PostmarkOptions
	_ = viper.BindEnv("postmarkOptions.serverToken", "POSTMARK_SERVER_TOKEN")
	_ = viper.BindEnv("postmarkOptions.fromEmail", "POSTMARK_FROM_EMAIL")

	// GomailOptions (legacy SMTP fallback)
	_ = viper.BindEnv("gomailOptions.host", "MAIL_HOST")
	_ = viper.BindEnv("gomailOptions.port", "MAIL_PORT")
	_ = viper.BindEnv("gomailOptions.username", "MAIL_USERNAME")
	_ = viper.BindEnv("gomailOptions.password", "MAIL_PASSWORD")
	_ = viper.BindEnv("gomailOptions.sender", "MAIL_SENDER")

	// WebsocketOptions
	_ = viper.BindEnv("websocketOptions.skipTLSVerification", "WS_SKIP_TLS_VERIFICATION")
	_ = viper.BindEnv("websocketOptions.originPatterns", "WS_ORIGIN_PATTERNS")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.WrapIf(err, "failed to unmarshal config")
	}

	if cfg.GormOptions.DBName == "" {
		return nil, errors.New("DB_NAME environment variable is required and was not found")
	}

	if cfg.EchoHttpOptions.SessionSecret == "" {
		return nil, errors.New("SESSION_SECRET environment variable is required and was not found")
	}

	return cfg, nil
}
