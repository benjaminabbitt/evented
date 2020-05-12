package client

import (
	"context"
	eventedbusiness "github.com/benjaminabbitt/evented/proto/evented/business"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewBusinessClient(target string, log *zap.SugaredLogger) (client BasicBusinessClient, err error) {
	log.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infow("Connected", "target", target)
	return BasicBusinessClient{
		log,
		eventedbusiness.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	log *zap.SugaredLogger
	bl  eventedbusiness.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *eventedcore.ContextualCommand) (events *eventedcore.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
