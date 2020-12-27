package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	business "github.com/benjaminabbitt/evented/proto/evented/business/coordinator"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/business/query"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	coordinator "github.com/benjaminabbitt/evented/proto/evented/saga/coordinator"
	"github.com/benjaminabbitt/evented/proto/evented/saga/saga"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient saga.SagaClient, eventQueryClient eventedquery.EventQueryClient, otherCommandHandlerClient business.BusinessCoordinatorClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
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
	coordinator.UnimplementedSagaCoordinatorServer
	coordinator *universal.SagaCoordinator
	Log         *zap.SugaredLogger
	Tracer      *opentracing.Tracer
}

func (o *Coordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar(), *o.Tracer)

	coordinator.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
