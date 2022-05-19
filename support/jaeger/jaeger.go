package jaeger

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	zap2jaeger "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
	"io"
)

func SetupJaeger(service string, log *zap.SugaredLogger) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
		ServiceName: service,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(zap2jaeger.NewLogger(log.Desugar())))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func CloseJaeger(closer io.Closer, log *zap.SugaredLogger) {
	if err := closer.Close(); err != nil {
		log.Errorw("Error closing Jaeger", err)
	}
}
