package config

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/mailing"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/postmark"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
)

type Config struct {
	Environment              environment.Environment `mapstructure:"ENVIRONMENT"`
	FrontendURL              string                  `mapstructure:"FRONTEND_URL"`
	RequireEmailVerification bool                    `mapstructure:"REQUIRE_EMAIL_VERIFICATION"`
	GormOptions              database.Options        `mapstructure:"GORMOPTIONS"`
	EchoHttpOptions          echoweb.Options         `mapstructure:"ECHOHTTPOPTIONS"`
	PostmarkOptions          postmark.Options        `mapstructure:"POSTMARKOPTIONS"`
	GomailOptions            mailing.Options         `mapstructure:"GOMAILOPTIONS"`
	WebsocketOptions         websocket.Options       `mapstructure:"WEBSOCKETOPTIONS"`
	LoggerOptions            logger.Options          `mapstructure:"LOGGEROPTIONS"`
}
