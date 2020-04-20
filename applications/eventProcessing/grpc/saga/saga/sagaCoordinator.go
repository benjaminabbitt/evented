package saga

import (
	"fmt"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/sagaCoordinator"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcZap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewSagaTracker(client evented_saga.SagaClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) SagaCoordinator {
	return SagaCoordinator{
		processed:  processedClient,
		log:        log,
		sagaClient: client,
		domain:     domain,
	}
}

type SagaCoordinator struct {
	evented_saga_coordinator.UnimplementedSagaCoordinatorServer
	domain              string //Domain of the Source
	log                 *zap.SugaredLogger
	sagaClient          evented_saga.SagaClient
	otherCommandHandler eventedcore.CommandHandlerClient
	eventQueryClient    evented_query.EventQueryClient
	processed           *processed.Processed
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book domain %s does not match saga configured domain %s", eb.Cover.Domain, o.domain))
	}
	id, err := evented_proto.ProtoToUUID(eb.Cover.Root)
	last, err := o.processed.LastReceived(ctx, id)
	seq := eb.Pages[0].Sequence.(*eventedcore.EventPage_Num).Num
	if err != nil {
		//TODO
	}
	if last < seq {
		evtStream, err := o.eventQueryClient.GetEvents(ctx, &evented_query.Query{
			Domain:     eb.Cover.Domain,
			Root:       eb.Cover.Root,
			LowerBound: seq,
		})
		if err != nil {
			o.log.Error(err)
		}
		for {
			event, err := evtStream.Recv()
			if err != nil {
				o.log.Error(err)
			}
			_, err = o.sagaClient.Handle(ctx, event)
			if err != nil {
				o.log.Error(err)
			} else {
				o.markProcessed(ctx, event)
			}
		}
	}

	reb, err := o.sagaClient.HandleSync(ctx, eb)
	if err != nil {
		o.log.Error(err)
	}
	o.markProcessed(ctx, eb)
	commandHandlerResponse, err := o.otherCommandHandler.Record(ctx, reb)
	commandHandlerResponse.Books = append(commandHandlerResponse.Books, reb)
	return commandHandlerResponse, err
}

func (o *SagaCoordinator) markProcessed(ctx context.Context, event *eventedcore.EventBook) {
	id, err := evented_proto.ProtoToUUID(event.Cover.Root)
	for _, page := range event.Pages {
		err = o.processed.Received(ctx, id, page.Sequence.(*eventedcore.EventPage_Num).Num)
		if err != nil {
			o.log.Error(err)
		}
	}
}

func (o *SagaCoordinator) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)

	grpcServer := grpcZap.GenerateConfiguredServer(o.log.Desugar())

	evented_saga_coordinator.RegisterSagaCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
