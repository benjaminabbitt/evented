package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"google.golang.org/grpc"
)

type GRPCSagaClient struct {
	client evented_saga.SagaClient
}

func (client GRPCSagaClient) SendSync(evts evented_core.EventBook) (responseEvents *evented_core.EventBook, err error) {
	return client.client.HandleSync(context.Background(), &evts)
}

func NewGRPCSagaClient() GRPCSagaClient {
	client := evented_saga.NewSagaClient(&grpc.ClientConn{})
	return GRPCSagaClient{client: client}
}
