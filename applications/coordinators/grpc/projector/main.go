package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/grpc/projector/projector"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	support.LogStartup(log, "GRPC Projector Coordinator Startup")

	config := configuration.Configuration{}
	config.Initialize(log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer jaeger.CloseJaeger(closer, log)

	target := config.ProjectorURL()
	log.Infow("Attempting to connect to Projector", "url", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	projectorClient := evented.NewProjectorClient(conn)

	healthClient := grpc_health_v1.NewHealthClient(conn)
	req := &grpc_health_v1.HealthCheckRequest{Service: "evented-sample-projector"}
	resp, err := healthClient.Check(context.Background(), req)
	log.Infow("Projector Status", "Health Check", resp)

	processedClient := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), config.CollectionName(), log)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)

	domain := config.Domain()

	lis, err := support.OpenPort(config.Port(), log)
	if err != nil {
		log.Error(err)
	}
	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	server := projector.NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, domain, log, &tracer)
	evented.RegisterProjectorCoordinatorServer(rpc, server)
	rpc.Serve(lis)
}
