package grpc

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/event/projector/configuration"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func ListenGRPC(log *zap.SugaredLogger, config *configuration.Configuration, tracer opentracing.Tracer) {
	target := config.Projector.Url
	log.Infow("Attempting to connect to Projector", "url", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	projectorClient := evented.NewProjectorClient(conn)

	healthClient := grpc_health_v1.NewHealthClient(conn)
	req := &grpc_health_v1.HealthCheckRequest{Service: "evented-sample-sample-projector"}
	resp, err := healthClient.Check(context.Background(), req)
	log.Infow("Projector Status", "Health Check", resp)

	processedClient := processed.NewProcessedClient(config.Database.Mongodb.Url, config.Database.Mongodb.Name, config.Database.Mongodb.Collection, log)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandler.Url, log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)

	domain := config.Domain

	lis, err := support.OpenPort(config.Port, log)
	if err != nil {
		log.Error(err)
	}
	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	server := NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, domain, log, &tracer)
	evented.RegisterProjectorCoordinatorServer(rpc, server)
	rpc.Serve(lis)
}
