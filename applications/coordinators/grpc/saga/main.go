package main

import (
	"github.com/benjaminabbitt/evented/applications/eventProcessing/grpc/saga/saga"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*
GRPC Server that receives Event messages and forwards them to a Sync Saga.
Fetches missing events from the event query server, if applicable.
Parses the result of the sync saga, updating last processed event in storage.
Sends the saga generated events to the other command handler
Returns the result of the sync saga.
*/
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

	target := viper.GetString("saga.url")
	sagaConn := grpcWithInterceptors.GenerateConfiguredConn(target, log)
	log.Infof("Connected to remote %s", target)
	if err != nil {
		log.Error(err)
	}
	sagaClient := evented_saga.NewSagaClient(sagaConn)

	ochUrl := viper.GetString("otherCommandHandler.url")
	otherCommandConn := grpcWithInterceptors.GenerateConfiguredConn(ochUrl, log)
	otherCommandHandler := evented_core.NewCommandHandlerClient(otherCommandConn)

	p := processed.NewProcessedClient(viper.GetString("database.url"), viper.GetString("database.name"), log)

	domain := viper.GetString("domain")

	server := saga.NewSagaCoordinator(sagaClient, otherCommandHandler, p, domain, log)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Saga Proxy Server...", "port", port)
	server.Listen(port)
}
