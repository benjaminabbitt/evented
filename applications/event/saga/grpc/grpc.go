package grpc

import (
	"github.com/benjaminabbitt/evented/applications/event/saga/configuration"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const NAME = "GRPC"

func ListenGRPC(log *zap.SugaredLogger, config *configuration.Config, tracer opentracing.Tracer) {
	sagaURL := config.Saga.Url
	sagaConn := grpcWithInterceptors.GenerateConfiguredConn(sagaURL, log, tracer)
	log.Infof("Connected to remote %s", sagaURL)
	sagaClient := evented.NewSagaClient(sagaConn)

	ochUrls := config.OtherCommandHandlers
	var ochConnections []evented.BusinessCoordinatorClient
	for _, ochUrl := range ochUrls {
		otherCommandConn := grpcWithInterceptors.GenerateConfiguredConn(ochUrl.Url, log, tracer)
		otherCommandHandler := evented.NewBusinessCoordinatorClient(otherCommandConn)
		ochConnections = append(ochConnections, otherCommandHandler)
	}

	p := processed.NewProcessedClient(config.Database.Mongodb.Url, config.Database.Mongodb.Name, config.Database.Mongodb.Collection, log)
	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandler.Url, log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)
	domain := config.Domain

	server := NewSagaCoordinator(sagaClient, eventQueryClient, ochConnections, p, domain, log, &tracer)

	port := config.Port
	log.Infow("Starting Saga Proxy Server...", "port", port)
	server.Listen(port)
}
