package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/integrationTest/eventSender/configuration"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/evented/sagaCoordinator"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("EventSender", log)

	log.Info("Starting...")
	target := config.EventHandlerURL()
	log.Info(target)

	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log)

	log.Infof("Connected to remote %s", target)

	sh := evented_saga_coordinator.NewSagaCoordinatorClient(conn)
	log.Info("Client Created...")

	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)

	var pages []*evented_core.EventPage
	for i := 0; i <= 1; i++ {
		pages = append(pages, framework.NewEventPage(uint32(i), false, nil))
	}
	eventBook := &evented_core.EventBook{
		Cover: &evented_core.Cover{
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
