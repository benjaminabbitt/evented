package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/integrationTest/eventSender/configuration"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	log.Info("Starting...")
	target := config.EventHandlerURL()
	log.Info(target)

	tracer, closer := jaeger.NewTracer(config.AppName(),
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)

	defer closer.Close()

	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)

	log.Infof("Connected to remote %s", target)

	sh := evented.NewSagaCoordinatorClient(conn)
	log.Info("Client Created...")

	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)

	var pages []*evented.EventPage
	for i := 0; i <= 1; i++ {
		pages = append(pages, framework.NewEventPage(uint32(i), false, nil))
	}
	eventBook := &evented.EventBook{
		Cover: &evented.Cover{
			Domain: config.Domain(),
			Root:   &protoId,
		},
		Pages: pages,
	}
	res, err := sh.HandleSync(context.Background(), eventBook)
	log.Info(res)
	if err != nil {

	}
}
