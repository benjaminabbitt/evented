package actx

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Actx struct {
	Log    *zap.SugaredLogger
	Tracer opentracing.Tracer
}
