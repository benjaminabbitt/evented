package actx

import (
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	"github.com/benjaminabbitt/evented/support/actx"
	"github.com/cenkalti/backoff/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type BasicCommandHandlerApplicationContext struct {
	actx.Actx
	Config        *configuration.Configuration
	RetryStrategy backoff.BackOff
}

func (o *BasicCommandHandlerApplicationContext) Log() *zap.SugaredLogger {
	return o.Actx.Log
}

func (o *BasicCommandHandlerApplicationContext) Tracer() opentracing.Tracer {
	return o.Actx.Tracer
}

func (o *BasicCommandHandlerApplicationContext) SetLog(logger *zap.SugaredLogger) {
	o.Actx.Log = logger
}

func (o *BasicCommandHandlerApplicationContext) SetTracer(tracer opentracing.Tracer) {
	o.Actx.Tracer = tracer
}
