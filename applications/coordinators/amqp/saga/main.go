package main

import (
	"context"
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/*
Dequeue from AMQP based message passing system,
*/
var log *zap.SugaredLogger

func main() {
	log = support.Log()
	defer log.Sync()

	config := Configuration{}
	config.Initialize(log)

	commandHandler := *makeCommandHandlerClient(config.CommandHandlerURL())

	ctx := context.Background()

	eh := makeEventHandlerClient(config)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	go func() {
		for {
			msg := <-decodedMessageChan
			reb, err := eh.Handle(ctx, msg.Book)
			if err != nil {
				log.Error(err)
				err = msg.Nack()
				if err != nil {
					log.Error(err)
				}
				continue
			}
			_, err = commandHandler.Record(ctx, reb)
			if err != nil {
				log.Error(err)
				err = msg.Nack()
				if err != nil {
					log.Error(err)
				}
				continue
			}
			err = msg.Ack()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	rabbitReceiver.ListenForever()
}

func makeRabbitReceiver(config Configuration) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
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

func makeEventHandlerClient(config Configuration) evented_saga.SagaClient {
	log.Info("Starting...")
	target := config.BusinessURL()
	log.Infow("Attempting to connect to business at", "address", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented_saga.NewSagaClient(conn)
	log.Info("Client Created...")
	return eventHandler
}

func makeCommandHandlerClient(target string) *evented_core.CommandHandlerClient {
	log.Infow("Attempting to connect to Command Handler at", "address", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	commandHandler := evented_core.NewCommandHandlerClient(conn)
	log.Info("Client Created...")
	return &commandHandler
}
