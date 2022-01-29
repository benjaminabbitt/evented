package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/business"

	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/saga"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/sagaCoordinator"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient saga.SagaClient, eventQueryClient evented.EventQueryClient, otherCommandHandlerClient business.BusinessCoordinatorClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
	universalCoordinator := &universal.Coordinator{
		Processed:        processedClient,
		EventQueryClient: eventQueryClient,
		Log:              log,
	}
	universalSagaCoordinator := &universal.SagaCoordinator{
		Coordinator:         universalCoordinator,
		Domain:              domain,
		SagaClient:          sagaClient,
		OtherCommandHandler: otherCommandHandlerClient,
	}
	return Coordinator{
		coordinator: universalSagaCoordinator,
		Log:         log,
		Tracer:      tracer,
	}
}

type Coordinator struct {
	sagaCoordinator.UnimplementedSagaCoordinatorServer
	coordinator *universal.SagaCoordinator
	Log         *zap.SugaredLogger
	Tracer      *opentracing.Tracer
}

func (o *Coordinator) HandleSync(ctx context.Context, eb *evented.EventBook) (*evented.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar(), *o.Tracer)

	sagaCoordinator.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
