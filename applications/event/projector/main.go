package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/event/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/event/projector/grpc/projector"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
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

	config := configuration.Configuration{}
	config.Initialize(log)

	ctx := context.Background()

	projectorClient := makeProjectorClient(config)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

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
		decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)
		listenRabbit(decodedMessageChan, rabbitReceiver, projectorCoordinator)
	} else if config.Transport.Kind == GRPC {
		//TODO: Unify approaches/contract here
		listenGRPC(ctx, &config, tracer)
	}

}

func listenGRPC(ctx context.Context, config *configuration.Configuration, tracer opentracing.Tracer) {
	target := config.Projector.Url
	log.Infow("Attempting to connect to Projector", "url", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	projectorClient := evented.NewProjectorClient(conn)

	healthClient := grpc_health_v1.NewHealthClient(conn)
	req := &grpc_health_v1.HealthCheckRequest{Service: "evented-sample-sample-projector"}
	resp, err := healthClient.Check(ctx, req)
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
	server := projector.NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, domain, log, &tracer)
	evented.RegisterProjectorCoordinatorServer(rpc, server)
	rpc.Serve(lis)
}

func listenRabbit(decodedMessageChan chan receiver.AMQPDecodedMessage, rabbitReceiver receiver.AMQPReceiver, coordinator *coordinator.ProjectorCoordinator) {
	for {
		msg := <-decodedMessageChan
		err := coordinator.Handle(context.Background(), msg.Book)
		if err == nil {
			err := msg.Ack()
			if err != nil {
				log.Error(err)
			}
		} else {
			err = msg.Nack()
			if err != nil {
				log.Error(err)
			}
		}
	}
	rabbitReceiver.ListenForever()
}

func makeRabbitReceiver(config configuration.Configuration) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
	outChan := make(chan receiver.AMQPDecodedMessage)
	receiverInstance := receiver.AMQPReceiver{
		SourceURL:         config.Transport.AMQP.Url,
		SourceExhangeName: config.Transport.AMQP.Exchange,
		SourceQueueName:   config.Transport.AMQP.Queue,
		Log:               log,
		OutputChannel:     outChan,
	}
	err := receiverInstance.Connect()
	if err != nil {
		log.Error(err)
	}
	log.Infow("Created RabbitMQ Receiver", "url", receiverInstance.SourceURL, "queue", receiverInstance.SourceQueueName)
	return outChan, receiverInstance
}

func makeProjectorClient(config configuration.Configuration) evented.ProjectorClient {
	log.Info("Starting...")
	target := config.Projector.Url
	log.Infow("Attempting to connect to projector at", "address", target)
	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented.NewProjectorClient(conn)
	log.Info("Client Created...")
	return eventHandler
}
