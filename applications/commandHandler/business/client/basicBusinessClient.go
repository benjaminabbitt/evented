package client

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBusinessClient(target string, log *zap.SugaredLogger) (client BasicBusinessClient, err error) {
	log.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to Business Logic: %v", err)
	}
	log.Infow("Connected", "target", target)
	return BasicBusinessClient{
		log,
		evented.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	log *zap.SugaredLogger
	bl  evented.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *evented.ContextualCommand) (events *evented.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
