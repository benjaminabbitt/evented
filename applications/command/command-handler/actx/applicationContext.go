package actx

import (
	"github.com/benjaminabbitt/evented/support/actx"
	"github.com/cenkalti/backoff/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BasicCommandHandlerApplicationContext struct {
	actx.Actx
	Config        *viper.Viper
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
