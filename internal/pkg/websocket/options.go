package websocket

import "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"

type Options struct {
	SkipTLSVerification bool     `mapstructure:"skipTLSVerification"`
	OriginPatterns      []string `mapstructure:"originPatterns"`
}

func ProvideConfig() (*Options, error) {
	return config.BindConfigKey[*Options]("websocketOptions")
}
