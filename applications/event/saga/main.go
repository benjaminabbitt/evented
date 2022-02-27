package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/event/saga/configuration"
	"github.com/benjaminabbitt/evented/applications/event/saga/grpc"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
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
	config.Initialize(log)

	log.Info("Configuration", config.Transport.Rabbitmq.Exchange)

	ctx := context.Background()

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	sagaClient := makeSagaClient(config, tracer)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandler.Url, log, tracer)
	eventQueryClient := evented.NewEventQueryClient(qhConn)

	ochConn := grpcWithInterceptors.GenerateConfiguredConn(config.OtherCommandHandlers[0].Url, log, tracer)
	otherCommandHandlerClient := evented.NewBusinessCoordinatorClient(ochConn)

	processedClient := processed.NewProcessedClient(config.Database.Mongodb.Url, config.Database.Mongodb.Name, config.Database.Mongodb.Collection, log)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	sagaCoordinator := coordinator.SagaCoordinator{
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

	go func() {
		for {
			msg := <-decodedMessageChan
			err := sagaCoordinator.Handle(ctx, msg.Book)
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
	}()
	rabbitReceiver.ListenForever()
}

func makeRabbitReceiver(config *configuration.Config) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
	outChan := make(chan receiver.AMQPDecodedMessage)
	receiverInstance := receiver.AMQPReceiver{
		SourceURL:         config.Transport.Rabbitmq.Url,
		SourceExhangeName: config.Transport.Rabbitmq.Exchange,
		SourceQueueName:   config.Transport.Rabbitmq.Queue,
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

func makeGRPCReceiver(config configuration.Config, tracer opentracing.Tracer) {
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

	server := grpc.NewSagaCoordinator(sagaClient, eventQueryClient, ochConnections, p, domain, log, &tracer)

	port := config.Port
	log.Infow("Starting Saga Proxy Server...", "port", port)
	server.Listen(port)
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
