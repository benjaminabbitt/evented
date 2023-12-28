package grpc

import (
	"context"
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient evented2.SagaClient, eventQueryClient evented2.EventQueryClient, otherCommandHandlerClient []evented2.BusinessCoordinatorClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
	universalCoordinator := &coordinator.Coordinator{
		Processed:        processedClient,
		EventQueryClient: eventQueryClient,
		Log:              log,
	}
	universalSagaCoordinator := &coordinator.SagaCoordinator{
		Coordinator:         universalCoordinator,
		Domain:              domain,
		SagaClient:          sagaClient,
		OtherCommandHandler: otherCommandHandlerClient[0],
	}
	return Coordinator{
		coordinator: universalSagaCoordinator,
		Log:         log,
		Tracer:      tracer,
	}
}

type Coordinator struct {
	evented2.UnimplementedSagaCoordinatorServer
	coordinator *coordinator.SagaCoordinator
	Log         *zap.SugaredLogger
	Tracer      *opentracing.Tracer
}

func (o *Coordinator) HandleSync(ctx context.Context, eb *evented2.EventBook) (*evented2.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar(), *o.Tracer)

	evented2.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
