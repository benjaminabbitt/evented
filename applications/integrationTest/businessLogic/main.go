package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
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
	server := businessLogic.NewMockBusinessLogic(log, errh)

	port := uint16(viper.GetUint("port"))
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port)
}


