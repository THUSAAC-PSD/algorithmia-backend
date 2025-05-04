package config

import (
	"fmt"
	"net/url"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config/environment"
)

type Config struct {
	EchoHTTPOptions EchoHTTPOptions `mapstructure:"echoHTTPOptions"`
}

func NewAppConfig(env environment.Environment) (*Config, error) {
	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type EchoHTTPOptions struct {
	Port                string   `mapstructure:"port"                validate:"required" env:"Port"`
	Development         bool     `mapstructure:"development"                             env:"Development"`
	BasePath            string   `mapstructure:"basePath"            validate:"required" env:"BasePath"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"                     env:"DebugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
	Timeout             int      `mapstructure:"timeout"                                 env:"Timeout"`
	Host                string   `mapstructure:"host"                                    env:"Host"`
	Name                string   `mapstructure:"name"                                    env:"Name"`
	SessionSecret       string   `mapstructure:"sessionSecret"                           env:"SessionSecret"`
}

func (c *EchoHTTPOptions) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c *EchoHTTPOptions) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}
