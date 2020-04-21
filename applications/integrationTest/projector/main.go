package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/projector/projector"
	"github.com/benjaminabbitt/evented/support"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

/*
Placeholder business logic -- used for Projector integration tests
*/
func main() {
	log = support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	log.Infow("Flags: ", "name", name, "configPath", configPath)

	_ = support.SetupConfig(name, configPath, flag.CommandLine)

	server := projector.NewPlaceholderProjectorLogic(log)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Projector Server...", "port", port)
	server.Listen(port)
}
