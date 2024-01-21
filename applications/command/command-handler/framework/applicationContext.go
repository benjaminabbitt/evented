package framework

import (
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	"github.com/benjaminabbitt/evented/support"
	"github.com/cenkalti/backoff/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type BasicCommandHandlerApplicationContext struct {
	support.BasicApplicationContext
	Tracer opentracing.Tracer
	Config *configuration.Configuration
}

func (o *BasicCommandHandlerApplicationContext) RetryStrategy() backoff.BackOff {
	return o.BasicApplicationContext.RetryStrategy
}

func (o *BasicCommandHandlerApplicationContext) Log() *zap.SugaredLogger {
	return o.BasicApplicationContext.Log
}

func (o *BasicCommandHandlerApplicationContext) GetConfig() *configuration.Configuration {
	return o.Config
}
