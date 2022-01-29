package main

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/projector"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
)

/*
GRPC Server that receives Event messages and forwards them to a Sync evented.
Fetches missing events from the event query server, if applicable.
Parses the result of the sync projector, updating last processed event in storage.
Returns the result of the sync evented.
*/
func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	target := config.ProjectorURL()
	log.Infow("Attempting to connect to Projector", "url", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	projectorClient := evented.NewProjectorClient(conn)

	processedClient := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)

	domain := config.Domain()

	server := projector.NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, domain, log, &tracer)

	port := config.Port()
	log.Infow("Starting Projector Proxy Server...", "port", port)
	server.Listen(port)
}
