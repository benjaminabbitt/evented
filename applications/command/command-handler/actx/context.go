package actx

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/cenkalti/backoff/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type BasicApplicationContext struct {
	support.BasicApplicationContext
	Tracer opentracing.Tracer
}

func (actx *BasicApplicationContext) Log() *zap.SugaredLogger {
	return actx.BasicApplicationContext.Log
}

func (actx *BasicApplicationContext) RetryStrategy() backoff.BackOff {
	return actx.BasicApplicationContext.RetryStrategy
}

type ApplicationContext interface {
	Log() *zap.SugaredLogger
	RetryStrategy() backoff.BackOff
}
