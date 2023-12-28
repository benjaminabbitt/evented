package actx

import (
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	"github.com/benjaminabbitt/evented/support"
	"github.com/cenkalti/backoff/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type ApplicationContext struct {
	support.BasicApplicationContext
	Tracer        opentracing.Tracer
	Configuration *configuration.Configuration
}

func (actx *ApplicationContext) Log() *zap.SugaredLogger {
	return actx.BasicApplicationContext.Log
}

func (actx *ApplicationContext) RetryStrategy() backoff.BackOff {
	return actx.BasicApplicationContext.RetryStrategy
}
