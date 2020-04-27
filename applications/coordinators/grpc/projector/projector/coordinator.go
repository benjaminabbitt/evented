package projector

import (
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
)

func NewProjectorCoordinator(client evented_projector.ProjectorClient, eventQueryClient evented_query.EventQueryClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger) ProjectorCoordinator {
	universalCoordinator := &universal.Coordinator{
		Processed:        processedClient,
		EventQueryClient: eventQueryClient,
		Log:              log,
	}
	universalProjectCoordinator := &universal.ProjectorCoordinator{
		Coordinator:      universalCoordinator,
		Domain:           domain,
		ProjectorClient:  client,
		Processed:        processedClient,
		EventQueryClient: eventQueryClient,
		Log:              log,
	}
	return ProjectorCoordinator{
		log:         log,
		Coordinator: universalProjectCoordinator,
	}
}

type ProjectorCoordinator struct {
	evented_projector_coordinator.UnimplementedProjectorCoordinatorServer
	Coordinator *universal.ProjectorCoordinator
	log         *zap.SugaredLogger
}

func (o *ProjectorCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.Projection, error) {
	return o.Coordinator.HandleSync(ctx, eb)
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
