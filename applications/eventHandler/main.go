package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/eventHandler/rabbitmq"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async/evented_amqp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/*
Transceiver.  Dequeue from event passing system and translate to GRPC calls
 */
var log *zap.SugaredLogger
var errh *evented.ErrLogger


func main() {
	log, errh = support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	errh.LogIfErr(err, "Error configuring application.")

	ch := makeEventHandlerClient(err)

	receiver := makeRabbitReceiver(ch)
	receiver.Listen()
}

func makeRabbitReceiver(ch evented_eventHandler.EventHandlerClient) rabbitmq.RabbitMQReceiver {
	receiver :=  rabbitmq.RabbitMQReceiver{
		SourceURL: viper.GetString("transport.source.amqp.url"),
		SourceExhangeName: viper.GetString("transport.source.amqp.exchange"),
		SourceQueueName:   viper.GetString("transport.source.amqp.queue"),
		Sender: evented_amqp.NewAMQPClient(
			viper.GetString("transport.target.amqp.url"),
			viper.GetString("transport.target.amqp.exchange"),
			log,
			errh,
		),
		Log:          log,
		Errh:         errh,
		EventHandler: ch,
	}
	log.Infow("Created RabbitMQ Receiver", "url", receiver.SourceURL, "queue", receiver.SourceQueueName)
	return receiver
}

func makeEventHandlerClient(err error) evented_eventHandler.EventHandlerClient {
	log.Info("Starting...")
	target := viper.GetString("business.address")
	log.Info(target)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
	log.Info(fmt.Sprintf("Connected to remote %s", target))
	errh.LogIfErr(err, fmt.Sprintf("Error dialing %s", target))
	ch := evented_eventHandler.NewEventHandlerClient(conn)
	log.Info("Client Created...")
	return ch
}


