package grpc

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient evented.SagaClient, eventQueryClient evented.EventQueryClient, otherCommandHandlerClient []evented.BusinessCoordinatorClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
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
	evented.UnimplementedSagaCoordinatorServer
	coordinator *coordinator.SagaCoordinator
	Log         *zap.SugaredLogger
	Tracer      *opentracing.Tracer
}

func (o *Coordinator) HandleSync(ctx context.Context, eb *evented.EventBook) (*evented.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar(), *o.Tracer)

	evented.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
