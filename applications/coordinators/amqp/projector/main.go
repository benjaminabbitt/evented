package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/amqp/projector/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
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
	config.Initialize("amqpEventCoordinator", log)

	ctx := context.Background()

	projectorClient := makeProjectorClient(config)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log)
	eventQueryClient := evented_query.NewEventQueryClient(qhConn)

	processedClient := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	sagaCoordinator := universal.NewProjectorCoordinator(projectorClient, eventQueryClient, processedClient, config.Domain(), log)

	go func() {
		for {
			msg := <-decodedMessageChan
			err := sagaCoordinator.Handle(ctx, msg.Book)
			if err != nil {
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
	log.Infow("Created RabbitMQ Receiver", "url", receiverInstance.SourceURL, "queue", receiverInstance.SourceQueueName)
	return outChan, receiverInstance
}

func makeProjectorClient(config configuration.Configuration) evented_projector.ProjectorClient {
	log.Info("Starting...")
	target := config.ProjectorURL()
	log.Infow("Attempting to connect to Projector at", "address", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented_projector.NewProjectorClient(conn)
	log.Info("Client Created...")
	return eventHandler
}

