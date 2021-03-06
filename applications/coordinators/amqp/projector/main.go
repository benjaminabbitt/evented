package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/amqp/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/business/query"
	"github.com/benjaminabbitt/evented/proto/evented/projector/projector"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"go.uber.org/zap"
)

/*
Dequeue from AMQP based message passing system,
*/
var log *zap.SugaredLogger

func main() {
	log = support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	ctx := context.Background()

	projectorClient := makeProjectorClient(config)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log, tracer)
	eventQueryClient := eventedquery.NewEventQueryClient(qhConn)

	processedClient := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	projectorCoordinator := universal.ProjectorCoordinator{
		Coordinator: &universal.Coordinator{
			Processed:        processedClient,
			EventQueryClient: eventQueryClient,
			Log:              log,
		},
		Domain:          config.Domain(),
		ProjectorClient: projectorClient,
		Log:             log,
	}

	go func() {
		for {
			msg := <-decodedMessageChan
			err := projectorCoordinator.Handle(ctx, msg.Book)
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

func makeRabbitReceiver(config configuration.Configuration) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
	outChan := make(chan receiver.AMQPDecodedMessage)
	receiverInstance := receiver.AMQPReceiver{
		SourceURL:         config.AMQPURL(),
		SourceExhangeName: config.AMQPExchange(),
		SourceQueueName:   config.AMQPQueue(),
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

func makeProjectorClient(config configuration.Configuration) projector.ProjectorClient {
	log.Info("Starting...")
	target := config.BusinessURL()
	log.Infow("Attempting to connect to business at", "address", target)
	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := projector.NewProjectorClient(conn)
	log.Info("Client Created...")
	return eventHandler
}
