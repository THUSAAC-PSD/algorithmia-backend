package gomail

import "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"

type Options struct {
	Host     string `mapstructure:"host"     validate:"required" env:"Host"`
	Port     int    `mapstructure:"port"     validate:"required" env:"Port"`
	Username string `mapstructure:"username" validate:"required" env:"Username"`
	Password string `mapstructure:"password" validate:"required" env:"Password"`
	Sender   string `mapstructure:"sender"   validate:"required" env:"Sender"`
}

func ProvideConfig() (*Options, error) {
	return config.BindConfigKey[*Options]("gomailOptions")
}
