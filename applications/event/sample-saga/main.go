package main

import (
	"github.com/benjaminabbitt/evented/applications/event/sample-saga/configuration"
	"github.com/benjaminabbitt/evented/applications/event/sample-saga/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/jaeger"
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
	config.Initialize(log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	server := saga.NewPlaceholderSagaLogic(log, &tracer)

	port := config.Port()
	log.Infow("Starting Saga Server...", "port", port)
	server.Listen(port)
}
