package projector

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/proto/evented/projector/coordinator"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client projectorCoordinator.ProjectorCoordinatorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *eventedcore.EventBook) (projection *eventedcore.Projection, err error) {
	err = backoff.Retry(func() error {
		projection, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return projection, err
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := projectorCoordinator.NewProjectorCoordinatorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
