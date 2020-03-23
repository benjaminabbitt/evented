package client

import (
	"context"
	evented_business "github.com/benjaminabbitt/evented/proto/business"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewBusinessClient(target string, log *zap.SugaredLogger)(client Client, err error){
	log.Infow("Setting up connection with Business Server..","target", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infow("Connected", "target", target)
	return Client{
		log,
		 evented_business.NewBusinessLogicClient(conn),
	}, nil
}

type Client struct {
	log *zap.SugaredLogger
	bl evented_business.BusinessLogicClient
}

func (client Client) Handle(command *evented_core.ContextualCommand)(events *evented_core.EventBook, err error){
	return client.bl.Handle(context.Background(), command)
}

