package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GRPCSagaClient struct {
	client evented.SagaCoordinatorClient
}

func (o GRPCSagaClient) HandleSync(ctx context.Context, evts *evented.EventBook, opts ...grpc.CallOption) (responses *evented.SynchronousProcessingResponse, err error) {
	err = backoff.Retry(func() error {
		responses, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return responses, err
}

func NewGRPCSagaClient(conn *grpc.ClientConn) GRPCSagaClient {
	client := evented.NewSagaCoordinatorClient(conn)
	return GRPCSagaClient{client: client}
}
