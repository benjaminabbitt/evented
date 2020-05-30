package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/configuration"
	"github.com/benjaminabbitt/evented/support"
	"github.com/uber/jaeger-client-go"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("businessLogic", log)

	tracer, closer := jaeger.NewTracer(*config.AppName,
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)

	defer closer.Close()

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port, tracer)
}
