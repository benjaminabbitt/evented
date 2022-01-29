package client

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/business"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewBusinessClient(target string, log *zap.SugaredLogger) (client BasicBusinessClient, err error) {
	log.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infow("Connected", "target", target)
	return BasicBusinessClient{
		log,
		business.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	log *zap.SugaredLogger
	bl  business.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *core.ContextualCommand) (events *core.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
