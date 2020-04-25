package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/universal"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"github.com/google/uuid"
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
	config.Initialize("amqpEventCoordinator", log)

	commandHandler := *makeCommandHandlerClient(config.CommandHandlerURL())
	coordinator := universal.Coordinator{
		Processed:        nil,
		EventQueryClient: nil,
		Log:              nil,
	}

	ctx := context.Background()

	sagaClient := makeSagaClient(config)

	decodedMessageChan, rabbitReceiver := makeRabbitReceiver(config)

	go func() {
		for {
			msg := <-decodedMessageChan
			coordinator.RepairSequencing(ctx, msg.Book, func(book *evented_core.EventBook) error {
				_, err := sagaClient.Handle(ctx, book)
				return err
			})
			reb, err := sagaClient.Handle(ctx, msg.Book)
			if err != nil {
				log.Error(err)
				err = msg.Nack()
				if err != nil {
					log.Error(err)
				}
				continue
			}
			recordedEvent, err := commandHandler.Record(ctx, reb)
			if err != nil {
				log.Error(err)
				err = msg.Nack()
				if err != nil {
					log.Error(err)
				}
				continue
			}

			coordinator.MarkProcessed(ctx,
				err = msg.Ack()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	rabbitReceiver.ListenForever()
}

func locateReturnedEventBook(root uuid.UUID, books []*evented_core.EventBook) *evented_core.EventBook {
	for _, book := range books {
		bookId, err := evented_proto.ProtoToUUID(book.Cover.Root)
		if err != nil {
			log.Error(err)
		}
		if root == bookId {
			return book
		}
	}
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

func makeSagaClient(config Configuration) evented_saga.SagaClient {
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
