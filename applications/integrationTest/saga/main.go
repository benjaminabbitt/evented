package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/saga/configuration"
	"github.com/benjaminabbitt/evented/applications/integrationTest/saga/saga"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

/*
Placeholder business logic -- used for Saga integration tests
*/
func main() {
	log = support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("saga", log)
	server := saga.NewPlaceholderSagaLogic(log)

	port := config.Port()
	log.Infow("Starting Saga Server...", "port", port)
	server.Listen(port)
}
