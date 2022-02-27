package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/event/projector/amqp"
	"github.com/benjaminabbitt/evented/applications/event/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/event/projector/grpc"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

/*
Dequeue from AMQP based message passing system,
*/
var log *zap.SugaredLogger

const RABBIT = "rabbitmq"
const GRPC = "grpc"

func main() {
	log = support.Log()
	defer log.Sync()
	support.LogStartup(log, "AMQP Projector Coordinator Startup")

	baseConfig := support.ConfigInit{}
	config := &configuration.Configuration{}
	config = baseConfig.Initialize(log, config).(*configuration.Configuration)

	tracer, closer := jaeger.SetupJaeger(config.Name, log)
	defer closer.Close()

	projectorClient := makeProjectorClient(config, tracer)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandler.Url, log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)
	healthClient := grpc_health_v1.NewHealthClient(qhConn)
	grpcHealth.HealthCheck(healthClient, config.QueryHandler.Name, log)
	processedClient := processed.NewProcessedClient(config.Database.Mongodb.Url, config.Database.Mongodb.Name, config.Database.Mongodb.Collection, log)

	projectorCoordinator := &coordinator.ProjectorCoordinator{
		Coordinator: &coordinator.Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
		Domain:          config.Domain,
		ProjectorClient: projectorClient,
		Log:             log,
	}

	//TODO: replace with future plugin framework if/when golang supports plugins in windows
	if config.Transport.Kind == RABBIT {
		decodedMessageChan, rabbitReceiver := amqp.MakeRabbitReceiver(log, config)
		go amqp.ListenRabbit(log, decodedMessageChan, rabbitReceiver, projectorCoordinator)
	} else if config.Transport.Kind == GRPC {
		//TODO: Unify approaches/contract here
		grpc.ListenGRPC(log, config, tracer)
	}

}

func makeProjectorClient(config *configuration.Configuration, tracer opentracing.Tracer) evented.ProjectorClient {
	log.Info("Starting...")
	target := config.Projector.Url
	log.Infow("Attempting to connect to projector at", "address", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented.NewProjectorClient(conn)
	log.Info("Client Created...")
	return eventHandler
}