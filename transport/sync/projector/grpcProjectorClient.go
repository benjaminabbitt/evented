package projector

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	eventedprojectorcoordinator "github.com/benjaminabbitt/evented/proto/evented/projectorCoordinator"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client eventedprojectorcoordinator.ProjectorCoordinatorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *eventedcore.EventBook) (projection *eventedcore.Projection, err error) {
	return o.client.HandleSync(ctx, evts)
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := eventedprojectorcoordinator.NewProjectorCoordinatorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
