package client

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/actx"
	evented "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewBusinessClient(actx *actx.ApplicationContext, target string) (client BasicBusinessClient, err error) {
	actx.Log().Infow("Setting up connection with Business Server..", "target", target)
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to Business Logic: %v", err)
	}
	actx.Log().Infow("Connected", "target", target)
	return BasicBusinessClient{
		actx,
		evented.NewBusinessLogicClient(conn),
	}, nil
}

type BasicBusinessClient struct {
	actx *actx.ApplicationContext
	bl   evented.BusinessLogicClient
}

func (client BasicBusinessClient) Handle(ctx context.Context, command *evented.ContextualCommand) (events *evented.EventBook, err error) {
	return client.bl.Handle(ctx, command)
}
