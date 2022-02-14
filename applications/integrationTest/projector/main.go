package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/integrationTest/projector/projector"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

/*
Placeholder business logic -- used for Projector integration tests
*/
func main() {
	log = support.Log()
	support.LogStartup(log, "Sample Projector")
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	server := projector.NewPlaceholderProjectorLogic(log, &tracer)

	port := config.Port()
	log.Infow("Starting Projector Server...", "port", port)
	server.Listen(port)
}
