package amqp

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/event/projector/configuration"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"go.uber.org/zap"
)

func ListenRabbit(log *zap.SugaredLogger, decodedMessageChan chan receiver.AMQPDecodedMessage, rabbitReceiver receiver.AMQPReceiver, coordinator *coordinator.ProjectorCoordinator) {
	for {
		msg := <-decodedMessageChan
		err := coordinator.Handle(context.Background(), msg.Book)
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
	rabbitReceiver.ListenForever()
}

func MakeRabbitReceiver(log *zap.SugaredLogger, config *configuration.Configuration) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
	outChan := make(chan receiver.AMQPDecodedMessage)
	receiverInstance := receiver.AMQPReceiver{
		SourceURL:         config.Transport.AMQP.Url,
		SourceExhangeName: config.Transport.AMQP.Exchange,
		SourceQueueName:   config.Transport.AMQP.Queue,
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
