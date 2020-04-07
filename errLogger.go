package evented

import "go.uber.org/zap"

type ErrLogger struct {
	Log *zap.SugaredLogger
}

func (el *ErrLogger) LogIfErr(err error, message string) {
	if err != nil {
		el.Log.Warnw(message, "err", err)
	}
}

func (el *ErrLogger) FailIfErr(err error, message string) {
	if err != nil {
		el.Log.Fatal(message, "err", err)
	}
}
