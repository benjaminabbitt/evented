package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/event/saga/amqp"
	"github.com/benjaminabbitt/evented/applications/event/saga/configuration"
	"github.com/benjaminabbitt/evented/applications/event/saga/grpc"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

/*
Dequeue from AMQP based message passing system,
*/
var log *zap.SugaredLogger

func main() {
	log = support.Log()
	defer log.Sync()

	support.LogStartup(log, "AMQP Saga Coordinator Startup")
	config := &configuration.Config{}
	support.Initialize(log, config)

	log.Info("Configuration", config.Transport.Rabbitmq.Exchange)

	tracer, closer := jaeger.SetupJaeger(config.Name, log)
	defer closer.Close()

	sagaClient := makeSagaClient(config, tracer)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandler.Url, log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)

	ochConn := grpcWithInterceptors.GenerateConfiguredConn(config.OtherCommandHandlers[0].Url, log, tracer)
	otherCommandHandlerClient := evented.NewBusinessCoordinatorClient(ochConn)

	processedClient := processed.NewProcessedClient(config.Database.Mongodb.Url, config.Database.Mongodb.Name, config.Database.Mongodb.Collection, log)

	sagaCoordinator := &coordinator.SagaCoordinator{
		Coordinator: &coordinator.Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
		Domain:              config.Domain,
		SagaClient:          sagaClient,
		OtherCommandHandler: otherCommandHandlerClient,
		Log:                 log,
	}

	if config.Transport.Kind == grpc.NAME {
		grpc.ListenGRPC(log, config, tracer)
	} else if config.Transport.Kind == amqp.NAME {
		decodedMessageChan, rabbitReceiver := amqp.MakeRabbitReceiver(log, config)
		go amqp.ListenRabbit(log, decodedMessageChan, rabbitReceiver, sagaCoordinator)
	}
}

func makeSagaClient(config *configuration.Config, tracer opentracing.Tracer) evented.SagaClient {
	log.Info("Starting...")
	target := config.Saga.Url
	log.Infow("Attempting to connect to business at", "address", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented.NewSagaClient(conn)
	log.Info("Client Created...")
	return eventHandler
}
