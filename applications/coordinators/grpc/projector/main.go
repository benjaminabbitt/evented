package main

import (
	"github.com/benjaminabbitt/evented/applications/eventProcessing/grpc/projector/projector"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

/*
GRPC Server that receives Event messages and forwards them to a Sync Projector.
Fetches missing events from the event query server, if applicable.
Parses the result of the sync projector, updating last processed event in storage.
Returns the result of the sync projector.
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

	target := viper.GetString("target.url")
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Infof("Connected to remote %s", target)
	if err != nil {
		log.Error(err)
	}
	client := evented_projector.NewProjectorClient(conn)

	p := processed.NewProcessedClient(viper.GetString("database.url"), viper.GetString("database.name"), log)

	domain := viper.GetString("domain")

	server := projector.NewProjectorCoordinator(client, p, domain, log)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Projector Proxy Server...", "port", port)
	server.Listen(port)
}
