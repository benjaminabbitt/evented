package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/evented/core"
	evented_projector_coordinator "github.com/benjaminabbitt/evented/proto/evented/projectorCoordinator"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client evented_projector_coordinator.ProjectorCoordinatorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error) {
	return o.client.HandleSync(ctx, evts)
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := evented_projector_coordinator.NewProjectorCoordinatorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
