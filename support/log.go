package support

import (
	"go.uber.org/zap"
)

func Log() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log := logger.Sugar()
	return log
}
