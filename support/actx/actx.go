package actx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Actx struct {
	context.Context
	Log    *zap.SugaredLogger
	Tracer opentracing.Tracer
	Config *viper.Viper
}
