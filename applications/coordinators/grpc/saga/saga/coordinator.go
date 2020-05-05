package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	evented_query "github.com/benjaminabbitt/evented/proto/evented/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/evented/saga"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/evented/sagaCoordinator"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient evented_saga.SagaClient, eventQueryClient evented_query.EventQueryClient, otherCommandHandlerClient eventedcore.CommandHandlerClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) SagaCoordinator {
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
	return SagaCoordinator{
		coordinator: universalSagaCoordinator,
		Log:         log,
	}
}

type SagaCoordinator struct {
	evented_saga_coordinator.UnimplementedSagaCoordinatorServer
	coordinator *universal.SagaCoordinator
	Log         *zap.SugaredLogger
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *SagaCoordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar())

	evented_saga_coordinator.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
