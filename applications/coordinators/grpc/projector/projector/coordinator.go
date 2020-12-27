package projector

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/business/query"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	projectorCoordinator "github.com/benjaminabbitt/evented/proto/evented/projector/coordinator"
	"github.com/benjaminabbitt/evented/proto/evented/projector/projector"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewProjectorCoordinator(client projector.ProjectorClient, eventQueryClient eventedquery.EventQueryClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
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
	return Coordinator{
		log:         log,
		Coordinator: universalProjectCoordinator,
		tracer:      tracer,
	}
}

type Coordinator struct {
	projectorCoordinator.UnimplementedProjectorCoordinatorServer
	Coordinator *universal.ProjectorCoordinator
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *Coordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.Projection, error) {
	return o.Coordinator.HandleSync(ctx, eb)
}

func (o *Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	projectorCoordinator.RegisterProjectorCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
