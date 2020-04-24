package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := Configuration{}
	config.Initialize(log)

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port)
}
