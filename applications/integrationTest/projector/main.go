package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/integrationTest/projector/projector"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

/*
Placeholder business logic -- used for Projector integration tests
*/
func main() {
	log = support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("projector", log)

	server := projector.NewPlaceholderProjectorLogic(log)

	port := config.Port()
	log.Infow("Starting Projector Server...", "port", port)
	server.Listen(port)
}
