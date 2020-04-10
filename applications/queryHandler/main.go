package main

import (
	businessLogic2 "github.com/benjaminabbitt/evented/applications/queryHandler/businessLogic"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	log, errh := support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	errh.LogIfErr(err, "Error configuring application.")
	repo, err := events.SetupEventRepo(log, errh)
	server := businessLogic2.NewEventQueryServer(repo, log, errh)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port)
}
