package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/coordinators/amqp/saga/saga"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/*
Transceiver.  Dequeue from event passing system and translate to GRPC calls
*/
var log *zap.SugaredLogger

func main() {
	log = support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	if err != nil {
		log.Error(err)
	}

	config := saga.Configuration{}

	commandHandler := *makeCommandHandlerClient(config.CommandHandlerURL())

	ctx := context.Background()

	eh := makeEventHandlerClient(config)

	ebChan, rec := makeRabbitReceiver(config)
	go func() {
		for {
			msg := <-ebChan
			reb, err := eh.Handle(ctx, msg.Book)
			if err != nil {
				log.Error(err)
				rec.NAck(msg.Tag)
				continue
			}
			_, err = commandHandler.Record(ctx, reb)
			if err != nil {
				log.Error(err)
				rec.NAck(msg.Tag)
				continue
			}
			err = rec.Ack(msg.Tag)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}()
	rec.ListenForever()
}

func makeRabbitReceiver(config saga.Configuration) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
	outChan := make(chan receiver.AMQPDecodedMessage)
	receiver := receiver.AMQPReceiver{
		SourceURL:         config.AMQPURL(),
		SourceExhangeName: config.AMQPExchange(),
		SourceQueueName:   config.AMQPQueue(),
		Log:               log,
		OutputChannel:     outChan,
	}
	log.Infow("Created RabbitMQ Receiver", "url", receiver.SourceURL, "queue", receiver.SourceQueueName)
	return outChan, receiver
}

func makeEventHandlerClient(config saga.Configuration) evented_saga.SagaClient {
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
