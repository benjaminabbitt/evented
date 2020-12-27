package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	sagaCoordinator "github.com/benjaminabbitt/evented/proto/evented/saga/coordinator"
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
)

type GRPCSagaClient struct {
	client sagaCoordinator.SagaCoordinatorClient
}

func (o GRPCSagaClient) HandleSync(ctx context.Context, evts *evented_core.EventBook) (responses *evented_core.SynchronousProcessingResponse, err error) {
	err = backoff.Retry(func() error {
		responses, err = o.client.HandleSync(ctx, evts)
		return err
	}, backoff.NewExponentialBackOff())
	return responses, err
}

func NewGRPCSagaClient(conn *grpc.ClientConn) GRPCSagaClient {
	client := sagaCoordinator.NewSagaCoordinatorClient(conn)
	return GRPCSagaClient{client: client}
}
