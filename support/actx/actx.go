package actx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Actx struct {
	context.Context
	Log    *zap.SugaredLogger
	Tracer opentracing.Tracer
}
