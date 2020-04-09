package client

import (
	"context"
	evented_business "github.com/benjaminabbitt/evented/proto/business"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewBusinessClient(target string, log *zap.SugaredLogger) (client BasicBusinessClient, err error) {
	log.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infow("Connected", "target", target)
	return BasicBusinessClient{
		log,
		evented_business.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	log *zap.SugaredLogger
	bl  evented_business.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *evented_core.ContextualCommand) (events *evented_core.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
