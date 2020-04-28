package main

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/projector"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
)

/*
GRPC Server that receives Event messages and forwards them to a Sync Projector.
Fetches missing events from the event query server, if applicable.
Parses the result of the sync projector, updating last processed event in storage.
Returns the result of the sync projector.
*/
func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("grpcProjectorCoordinator", log)

	target := config.ProjectorURL()
	log.Infow("Attempting to connect to Projector", "url", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log)
	projectorClient := evented_projector.NewProjectorClient(conn)

	processedClient := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log)
	eventQueryClient := evented_query.NewEventQueryClient(qhConn)

	domain := config.Domain()

	server := projector.NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, domain, log)

	port := config.Port()
	log.Infow("Starting Projector Proxy Server...", "port", port)
	server.Listen(port)
}