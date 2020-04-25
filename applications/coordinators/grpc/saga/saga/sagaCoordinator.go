package saga

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/sagaCoordinator"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewSagaCoordinator(sagaClient evented_saga.SagaClient, eventQueryClient evented_query.EventQueryClient, otherCommandHandlerClient eventedcore.CommandHandlerClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) SagaCoordinator {
	return SagaCoordinator{
		otherCommandHandler: otherCommandHandlerClient,
		sagaClient:          sagaClient,
		domain:              domain,
		Coordinator: universal.Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
	}
}

type SagaCoordinator struct {
	evented_saga_coordinator.UnimplementedSagaCoordinatorServer
	universal.Coordinator
	domain              string
	sagaClient          evented_saga.SagaClient
	otherCommandHandler eventedcore.CommandHandlerClient
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book domain %s does not match saga configured domain %s", eb.Cover.Domain, o.domain))
	}

	o.RepairSequencing(ctx, eb, func(eb *eventedcore.EventBook) error {
		_, err := o.sagaClient.Handle(ctx, eb)
		return err
	})

	reb, err := o.sagaClient.HandleSync(ctx, eb)
	if err != nil {
		o.Log.Error(err)
	}
	o.MarkProcessed(ctx, eb)
	commandHandlerResponse, err := o.otherCommandHandler.Record(ctx, reb)
	commandHandlerResponse.Books = append(commandHandlerResponse.Books, reb)
	return commandHandlerResponse, err
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
