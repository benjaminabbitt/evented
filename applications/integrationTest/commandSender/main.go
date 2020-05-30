package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/integrationTest/commandSender/configuration"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/google/uuid"
	"time"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("commandSender", log)

	log.Info("Starting...")
	target := config.CommandHandlerURL()
	log.Info(target)
	tracer, closer := jaeger.SetupJaeger("commandSender", log)
	defer closer.Close()

	span := tracer.StartSpan("test")
	time.Sleep(1 * time.Second)
	span.Finish()

	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	log.Infof("Connected to remote %s", target)
	ch := evented_core.NewCommandHandlerClient(conn)
	log.Info("Client Created...")
	id, _ := uuid.NewRandom()
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
