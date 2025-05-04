package echoweb

import (
	"fmt"
	"net/url"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
)

type Options struct {
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

func (c *Options) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (c *Options) BasePathAddress() string {
	path, err := url.JoinPath(c.Address(), c.BasePath)
	if err != nil {
		return ""
	}
	return path
}

func ProvideConfig() (*Options, error) {
	return config.BindConfigKey[*Options]("echoHttpOptions")
}
