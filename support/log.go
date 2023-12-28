package support

import (
	"go.uber.org/zap"
)

func Log() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment(zap.AddCaller())
	//logger, _ := zap.NewProduction(zap.AddCaller())
	log := logger.Sugar()
	return log
}
