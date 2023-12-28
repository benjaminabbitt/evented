package amqp

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/event/saga/configuration"
	"github.com/benjaminabbitt/evented/support/coordinator"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"go.uber.org/zap"
)

const NAME = "RABBIT"

func MakeRabbitReceiver(log *zap.SugaredLogger, config *configuration.Config) (chan receiver.AMQPDecodedMessage, receiver.AMQPReceiver) {
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

func ListenRabbit(log *zap.SugaredLogger, decodedMessageChan chan receiver.AMQPDecodedMessage, rabbitReceiver receiver.AMQPReceiver, coordinator *coordinator.SagaCoordinator) {
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
}
