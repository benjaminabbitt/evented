package main

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/saga/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/saga/saga"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
)

/*
GRPC Server that receives Event messages and forwards them to a Sync Saga.
Fetches missing events from the event query server, if applicable.
Parses the result of the sync saga, updating last processed event in storage.
Sends the saga generated events to the other command handler
Returns the result of the sync saga.
*/
func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("grpcSagaCoordinator", log)

	sagaURL := config.SagaURL()
	sagaConn := grpcWithInterceptors.GenerateConfiguredConn(sagaURL, log)
	log.Infof("Connected to remote %s", sagaURL)
	sagaClient := evented_saga.NewSagaClient(sagaConn)

	ochUrl := config.OtherCommandHandlerURL()
	otherCommandConn := grpcWithInterceptors.GenerateConfiguredConn(ochUrl, log)
	otherCommandHandler := evented_core.NewCommandHandlerClient(otherCommandConn)

	p := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)
	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log)
	eventQueryClient := evented_query.NewEventQueryClient(qhConn)
	domain := config.Domain()

	server := saga.NewSagaCoordinator(sagaClient, eventQueryClient, otherCommandHandler, p, domain, log)

	port := config.Port()
	log.Infow("Starting Saga Proxy Server...", "port", port)
	server.Listen(port)
}
