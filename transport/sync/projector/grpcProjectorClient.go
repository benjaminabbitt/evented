package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client evented.ProjectorCoordinatorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *evented.EventBook, opts ...grpc.CallOption) (projection *evented.Projection, err error) {
	err = backoff.Retry(func() error {
		projection, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return projection, err
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := evented.NewProjectorCoordinatorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
