package main

import (
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	commandHandlers := make(map[string]evented_core.CommandHandlerClient)
	config := viper.Get("commandHandlers")
	log.Info(config)

	for key, element := range config.(map[string]interface{}) {
		log.Info(key)
		log.Info(element)
		url := element.(map[string]interface{})["url"]
		log.Info(url)
		commandHandlers[key] = *makeCommandHandlerClient(url.(string))
	}

	eh := makeEventHandlerClient()

	receiver := makeRabbitReceiver(*eh, commandHandlers)
	receiver.Listen()
}

func makeRabbitReceiver(
	eventHandler evented_eventHandler.EventHandlerClient,
	commandHandlers map[string]evented_core.CommandHandlerClient) receiver.AMQPReceiver {
	receiver := receiver.AMQPReceiver{
		SourceURL:         viper.GetString("transport.source.amqp.url"),
		SourceExhangeName: viper.GetString("transport.source.amqp.exchange"),
		SourceQueueName:   viper.GetString("transport.source.amqp.queue"),
		DestinationSink:   commandHandlers,
		Log:               log,
		EventHandler:      eventHandler,
	}
	log.Infow("Created RabbitMQ Receiver", "url", receiver.SourceURL, "queue", receiver.SourceQueueName)
	return receiver
}

func makeEventHandlerClient() *evented_eventHandler.EventHandlerClient {
	log.Info("Starting...")
	target := viper.GetString("business.address")
	log.Infow("Attempting to connect to business at", "address", target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error(err)
	}
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	eventHandler := evented_eventHandler.NewEventHandlerClient(conn)
	log.Info("Client Created...")
	return &eventHandler
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
