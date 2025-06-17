package config

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/mailing"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
)

type Config struct {
	Environment      environment.Environment `mapstructure:"ENVIRONMENT"`
	GormOptions      database.Options        `mapstructure:"GORMOPTIONS"`
	EchoHttpOptions  echoweb.Options         `mapstructure:"ECHOHTTPOPTIONS"`
	GomailOptions    mailing.Options         `mapstructure:"GOMAILOPTIONS"`
	WebsocketOptions websocket.Options       `mapstructure:"WEBSOCKETOPTIONS"`
	LoggerOptions    logger.Options          `mapstructure:"LOGGEROPTIONS"`
}
