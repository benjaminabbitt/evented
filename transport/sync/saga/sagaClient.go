package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"google.golang.org/grpc"
	"log"
)

type GRPCSagaClient struct {
	client evented_saga.SagaClient
}

func (client GRPCSagaClient) SendSync(evts evented_core.EventBook) (responseEvents *evented_core.EventBook, err error) {
	responseEvents, err = client.client.HandleSync(context.Background(), &evts)
	if err != nil {
		log.Fatal(err)
	}
	return responseEvents, nil
}

func NewGRPCSagaClient() GRPCSagaClient {
	client := evented_saga.NewSagaClient(&grpc.ClientConn{})
	return GRPCSagaClient{client: client}
}
