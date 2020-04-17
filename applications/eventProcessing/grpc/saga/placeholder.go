package main

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewSagaTracker(client evented_saga.SagaClient, processedClient *processed.Processed, log *zap.SugaredLogger) SagaTracker {
	return SagaTracker{
		processed: processedClient,
		log:       log,
		client:    client,
	}
}

type SagaTracker struct {
	evented_saga.UnimplementedSagaServer
	log       *zap.SugaredLogger
	client    evented_saga.SagaClient
	processed *processed.Processed
}

func (o *SagaTracker) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.EventBook, error) {
	reb, err := o.client.HandleSync(ctx, eb)
	uuid, err := evented_proto.ProtoToUUID(eb.Cover.Root)
	err = o.processed.Received(ctx, uuid, eb.Pages[0].Sequence.(*eventedcore.EventPage_Num).Num)
	if err != nil {
		o.log.Error(err)
	}
	return reb, err
}

func (o *SagaTracker) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpc.NewServer()

	evented_saga.RegisterSagaServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
