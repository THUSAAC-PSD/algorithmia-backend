package defaultlogger

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"
)

var l logger.Logger

func initLogger() {
	l = logger.NewZapLogger(
		&logger.LogOptions{CallerEnabled: false},
		constant.Dev,
	)
}

func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
