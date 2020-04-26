package universal

import (
	"fmt"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewSagaCoordinator(sagaClient evented_saga.SagaClient, eventQueryClient evented_query.EventQueryClient, otherCommandHandlerClient eventedcore.CommandHandlerClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) SagaCoordinator {
	return SagaCoordinator{
		OtherCommandHandler: otherCommandHandlerClient,
		SagaClient:          sagaClient,
		Domain:              domain,
		Coordinator: Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
	}
}

type SagaCoordinator struct {
	Coordinator
	Domain              string
	SagaClient          evented_saga.SagaClient
	OtherCommandHandler eventedcore.CommandHandlerClient
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match saga configured Domain %s", eb.Cover.Domain, o.Domain))
	}

	o.RepairSequencing(ctx, eb, func(eb *eventedcore.EventBook) error {
		_, err := o.SagaClient.Handle(ctx, eb)
		return err
	})

	reb, err := o.SagaClient.HandleSync(ctx, eb)
	if err != nil {
		o.Log.Error(err)
	}
	o.MarkProcessed(ctx, eb)
	commandHandlerResponse, err := o.OtherCommandHandler.Record(ctx, reb)
	commandHandlerResponse.Books = append(commandHandlerResponse.Books, reb)
	return commandHandlerResponse, err
}

func (o *SagaCoordinator) Handle(ctx context.Context, eb *eventedcore.EventBook) (err error) {
	_, err = o.HandleSync(ctx, eb)
	return err
}
