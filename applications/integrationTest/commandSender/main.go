package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/integrationTest/commandSender/configuration"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("commandSender", log)

	log.Info("Starting...")
	target := config.CommandHandlerURL()
	log.Info(target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infof("Connected to remote %s", target)
	if err != nil {
		log.Error(err)
	}
	ch := evented_core.NewCommandHandlerClient(conn)
	log.Info("Client Created...")
	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)

	for i := 0; i <= 1; i++ {
		pages := []*evented_core.CommandPage{&evented_core.CommandPage{
			Sequence:    uint32(i),
			Synchronous: false,
			Command:     nil,
		}}
		commandBook := &evented_core.CommandBook{
			Cover: &evented_core.Cover{
				Domain: config.Domain(),
				Root:   &protoId,
			},
			Pages: pages,
		}
		_, _ = ch.Handle(context.Background(), commandBook)
	}

	log.Info("Done!")
}
