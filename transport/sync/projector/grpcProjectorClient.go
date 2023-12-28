package projector

import (
	"context"
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GrpcProjector struct {
	client evented2.ProjectorCoordinatorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *evented2.EventBook, opts ...grpc.CallOption) (projection *evented2.Projection, err error) {
	err = backoff.Retry(func() error {
		projection, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return projection, err
}

func NewGRPCProjector(conn *grpc.ClientConn) GrpcProjector {
	client := evented2.NewProjectorCoordinatorClient(conn)
	return GrpcProjector{
		client: client,
	}
}
