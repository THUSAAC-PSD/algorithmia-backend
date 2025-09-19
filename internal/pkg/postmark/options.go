package postmark

type Options struct {
	ServerToken string `mapstructure:"serverToken" validate:"required" env:"ServerToken"`
	FromEmail   string `mapstructure:"fromEmail"   validate:"required" env:"FromEmail"`
}
