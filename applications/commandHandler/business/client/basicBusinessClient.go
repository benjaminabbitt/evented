package client

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewBusinessClient(target string, log *zap.SugaredLogger) (client BasicBusinessClient, err error) {
	log.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
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
