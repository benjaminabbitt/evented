package saga

import (
	"context"
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GRPCSagaClient struct {
	client evented2.SagaCoordinatorClient
}

func (o GRPCSagaClient) HandleSync(ctx context.Context, evts *evented2.EventBook, opts ...grpc.CallOption) (responses *evented2.SynchronousProcessingResponse, err error) {
	err = backoff.Retry(func() error {
		responses, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return responses, err
}

func NewGRPCSagaClient(conn *grpc.ClientConn) GRPCSagaClient {
	client := evented2.NewSagaCoordinatorClient(conn)
	return GRPCSagaClient{client: client}
}
