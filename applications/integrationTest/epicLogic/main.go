package main

import (
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/integrationTest/epicLogic/epicLogic"
	"github.com/benjaminabbitt/evented/support"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	log, errh = support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	log.Infow("Flags: ", "name", name, "configPath", configPath)

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	errh.LogIfErr(err, "Error configuring application.")

	errh = &evented.ErrLogger{log}

	server := epicLogic.NewMockEpicLogic(log, errh)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Epic Server...", "port", port)
	server.Listen(port)
}

