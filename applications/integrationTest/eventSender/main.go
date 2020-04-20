package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/sagaCoordinator"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcZap"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	var name = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	if err != nil {
		log.Error(err)
	}

	log.Info("Starting...")
	target := viper.GetString("eventHandlerURL")
	log.Info(target)

	conn := grpcZap.GenerateConfiguredConn(target, log)

	log.Infof("Connected to remote %s", target)
	if err != nil {
		log.Error(err)
		stat, _ := status.FromError(err)
		log.Error(stat)
	}

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
			Domain: viper.GetString("domain"),
			Root:   &protoId,
		},
		Pages: pages,
	}
	res, err := sh.HandleSync(context.Background(), eventBook)
	log.Info(res)
	if err != nil {

	}
}
