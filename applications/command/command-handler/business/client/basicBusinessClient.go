package client

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework/actx"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewBusinessClient(actx actx.CommandHandlerContext, target string) (client BasicBusinessClient, err error) {
	actx.Logger.Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to Business Logic: %v", err)
	}
	actx.Logger.Infow("Connected", "target", target)
	return BasicBusinessClient{
		actx,
		evented.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	actx actx.CommandHandlerContext
	bl   evented.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *evented.ContextualCommand) (events *evented.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
