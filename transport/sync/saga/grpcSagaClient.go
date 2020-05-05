package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/evented/sagaCoordinator"
	"google.golang.org/grpc"
)

type GRPCSagaClient struct {
	client evented_saga_coordinator.SagaCoordinatorClient
}

func (o GRPCSagaClient) HandleSync(ctx context.Context, evts *evented_core.EventBook) (responses *evented_core.SynchronousProcessingResponse, err error) {
	return o.client.HandleSync(ctx, evts)
}

func NewGRPCSagaClient(conn *grpc.ClientConn) GRPCSagaClient {
	client := evented_saga_coordinator.NewSagaCoordinatorClient(conn)
	return GRPCSagaClient{client: client}
}
