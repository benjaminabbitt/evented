package main

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	if err != nil {
		log.Error(err)
	}

	log.Info("Starting...")
	target := viper.GetString("commandHandlerURL")
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
				Domain: viper.GetString("domain"),
				Root:   &protoId,
			},
			Pages: pages,
		}
		_, _ = ch.Handle(context.Background(), commandBook)
	}

	log.Info("Done!")
}
