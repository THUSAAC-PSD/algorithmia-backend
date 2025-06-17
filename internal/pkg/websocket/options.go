package websocket

type Options struct {
	SkipTLSVerification bool     `mapstructure:"skipTLSVerification"`
	OriginPatterns      []string `mapstructure:"originPatterns"`
}
