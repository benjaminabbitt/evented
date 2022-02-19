package projector

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewProjectorCoordinator(client evented.ProjectorClient, eventQueryClient evented.EventQueryClient, processedClient *processed.Processed, domain string, log *zap.SugaredLogger, tracer *opentracing.Tracer) Coordinator {
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
	evented.UnimplementedProjectorCoordinatorServer
	Coordinator *universal.ProjectorCoordinator
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o Coordinator) HandleSync(ctx context.Context, eb *evented.EventBook) (*evented.Projection, error) {
	return o.Coordinator.HandleSync(ctx, eb)
}

func (o Coordinator) Listen(port uint) {
	lis := support.CreateListener(port, o.log)

	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	evented.RegisterProjectorCoordinatorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}

func (o Coordinator) Handle(ctx context.Context, eb *evented.EventBook) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, o.Coordinator.Handle(ctx, eb)
}
