package echoweb

type Options struct {
	Port          string `mapstructure:"port"          validate:"required" env:"Port"`
	SessionSecret string `mapstructure:"sessionSecret"                     env:"SessionSecret"`
}
