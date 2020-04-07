package support

import (
	"github.com/benjaminabbitt/evented"
	"go.uber.org/zap"
)

func Log() (*zap.SugaredLogger, *evented.ErrLogger) {
	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log := logger.Sugar()
	errh := &evented.ErrLogger{log}
	return log, errh
}
