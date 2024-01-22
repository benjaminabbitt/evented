package support

import (
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

type BasicApplicationContext struct {
	RetryStrategy backoff.BackOff
	Logger        *zap.SugaredLogger
	Domain        string
}

//func (b BasicApplicationContext) Logger() *zap.SugaredLogger {
//	return b.logger
//}
//
//func (b BasicApplicationContext) RetryStrategy() backoff.BackOff {
//	return b.retryStrategy
//}
//
//func (b BasicApplicationContext) Domain() string {
//	return b.domain
//}

type ApplicationContext interface {
	RetryStrategy() backoff.BackOff
	Logger() *zap.SugaredLogger
	Domain() string
}
