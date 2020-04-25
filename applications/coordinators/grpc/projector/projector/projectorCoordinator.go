package projector

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
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

func NewProjectorCoordinator(client evented_projector.ProjectorClient, eventQueryClient evented_query.EventQueryClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) ProjectorCoordinator {
	return ProjectorCoordinator{
		processed:       processedClient,
		log:             log,
		projectorClient: client,
		domain:          domain,
		Coordinator: universal.Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
	}
}

type ProjectorCoordinator struct {
	evented_projector_coordinator.UnimplementedProjectorCoordinatorServer
	universal.Coordinator
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
	o.RepairSequencing(ctx, eb, func(eb *eventedcore.EventBook) error {
		_, err := o.projectorClient.Handle(ctx, eb)
		return err
	})

	reb, err := o.projectorClient.HandleSync(ctx, eb)
	if err != nil {
		o.log.Error(err)
	}
	o.MarkProcessed(ctx, eb)
	return reb, err
}

func (o *ProjectorCoordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar())

	evented_projector_coordinator.RegisterProjectorCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
