package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client evented_projector.ProjectorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error) {
	return o.client.HandleSync(ctx, evts)
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := evented_projector.NewProjectorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
