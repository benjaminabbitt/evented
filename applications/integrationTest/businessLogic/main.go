package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/configuration"
	"github.com/benjaminabbitt/evented/support"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("businessLogic", log)

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port)
}
