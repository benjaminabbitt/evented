package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/amqp/configuration"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
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

	sagaClient := makeSagaClient(config)

	qhConn := grpcWithInterceptors.GenerateConfiguredConn(config.QueryHandlerURL(), log)
	eventQueryClient := evented_query.NewEventQueryClient(qhConn)

	ochConn := grpcWithInterceptors.GenerateConfiguredConn(config.CommandHandlerURL(), log)
	otherCommandHandlerClient := evented_core.NewCommandHandlerClient(ochConn)

	processed := processed.NewProcessedClient(config.DatabaseURL(), config.DatabaseName(), log)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	sagaCoordinator := universal.NewSagaCoordinator(sagaClient, eventQueryClient, otherCommandHandlerClient, processed, config.Domain(), log)

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

func makeSagaClient(config configuration.Configuration) evented_saga.SagaClient {
	log.Info("Starting...")
	target := config.BusinessURL()
	log.Infow("Attempting to connect to business at", "address", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented_saga.NewSagaClient(conn)
	log.Info("Client Created...")
	return eventHandler
}

func makeCommandHandlerClient(target string) *evented_core.CommandHandlerClient {
	log.Infow("Attempting to connect to Command Handler at", "address", target)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log)
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	commandHandler := evented_core.NewCommandHandlerClient(conn)
	log.Info("Client Created...")
	return &commandHandler
}
