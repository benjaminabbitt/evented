package support

import (
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

type BasicApplicationContext struct {
	RetryStrategy backoff.BackOff
	Log           *zap.SugaredLogger
}
