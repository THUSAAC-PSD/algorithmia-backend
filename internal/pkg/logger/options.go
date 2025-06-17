package logger

type Options struct {
	LogLevel      string `mapstructure:"level"`
	CallerEnabled bool   `mapstructure:"callerEnabled"`
}
