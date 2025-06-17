package defaultlogger

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/environment"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
)

var l logger.Logger

func initLogger() {
	l = logger.NewZapLogger(
		&logger.Options{CallerEnabled: false},
		&environment.Development,
	)
}

func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
