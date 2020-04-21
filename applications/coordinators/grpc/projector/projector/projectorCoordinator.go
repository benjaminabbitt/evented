package projector

import (
	"fmt"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	evented_projector_coordinator "github.com/benjaminabbitt/evented/proto/projectorCoordinator"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewProjectorCoordinator(client evented_projector.ProjectorClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) ProjectorCoordinator {
	return ProjectorCoordinator{
		processed:       processedClient,
		log:             log,
		projectorClient: client,
		domain:          domain,
	}
}

type ProjectorCoordinator struct {
	evented_projector_coordinator.UnimplementedProjectorCoordinatorServer
	domain           string //Domain of the Source
	log              *zap.SugaredLogger
	projectorClient  evented_projector.ProjectorClient
	eventQueryClient evented_query.EventQueryClient
	processed        *processed.Processed
}

func (o *ProjectorCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.Projection, error) {
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
			_, err = o.projectorClient.Handle(ctx, event)
			if err != nil {
				o.log.Error(err)
			} else {
				o.markProcessed(ctx, event)
			}
		}
	}

	reb, err := o.projectorClient.HandleSync(ctx, eb)
	if err != nil {
		o.log.Error(err)
	}
	o.markProcessed(ctx, eb)
	return reb, err
}

func (o *ProjectorCoordinator) markProcessed(ctx context.Context, event *eventedcore.EventBook) {
	id, err := evented_proto.ProtoToUUID(event.Cover.Root)
	for _, page := range event.Pages {
		err = o.processed.Received(ctx, id, page.Sequence.(*eventedcore.EventPage_Num).Num)
		if err != nil {
			o.log.Error(err)
		}
	}
}

func (o *ProjectorCoordinator) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar())

	evented_projector_coordinator.RegisterProjectorCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
