package main

import (
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/saga/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/saga/saga"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
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

	support.LogStartup(log, "GRPC Saga Coordinator Startup")

	config := configuration.Configuration{}
	config.Initialize(log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	sagaURL := config.SagaURL()
	sagaConn := grpcWithInterceptors.GenerateConfiguredConn(sagaURL, log, tracer)
	log.Infof("Connected to remote %s", sagaURL)
	sagaClient := evented.NewSagaClient(sagaConn)

	ochUrl := config.OtherCommandHandlerURL()
	otherCommandConn := grpcWithInterceptors.GenerateConfiguredConn(ochUrl, log, tracer)
	//TODO: fix
	otherCommandHandler := evented.NewBusinessCoordinatorClient(otherCommandConn)

	p := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)
	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)
	domain := config.Domain()

	server := saga.NewSagaCoordinator(sagaClient, eventQueryClient, otherCommandHandler, p, domain, log, &tracer)

	port := config.Port()
	log.Infow("Starting Saga Proxy Server...", "port", port)
	server.Listen(port)
}
