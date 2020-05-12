package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/query"
	eventedsaga "github.com/benjaminabbitt/evented/proto/evented/saga"
	eventedsagacoordinator "github.com/benjaminabbitt/evented/proto/evented/sagaCoordinator"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"go.uber.org/zap"
)

func NewSagaCoordinator(sagaClient eventedsaga.SagaClient, eventQueryClient eventedquery.EventQueryClient, otherCommandHandlerClient eventedcore.CommandHandlerClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) SagaCoordinator {
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
	eventedsagacoordinator.UnimplementedSagaCoordinatorServer
	coordinator *universal.SagaCoordinator
	Log         *zap.SugaredLogger
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	return o.coordinator.HandleSync(ctx, eb)
}

func (o *SagaCoordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.Log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.Log.Desugar())

	eventedsagacoordinator.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.Log.Error(err)
	}
}
