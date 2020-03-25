package rabbitmq

import (
	"context"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/benjaminabbitt/evented/transport/async/evented_amqp"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)


type RabbitMQReceiver struct {
	SourceURL			string
	SourceExhangeName string
	SourceQueueName   string
	Sender            *evented_amqp.Client
	Errh              *evented.ErrLogger
	Log               *zap.SugaredLogger
	EventHandler      evented_eventHandler.EventHandlerClient
}

func (mq *RabbitMQReceiver) Listen(){
	viper.GetString("")

	conn, err := amqp.Dial(mq.SourceURL)
	mq.Errh.LogIfErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mq.Errh.LogIfErr(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		mq.SourceExhangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	mq.Errh.LogIfErr(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		mq.SourceQueueName,
		false,   // durable
		false,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	mq.Errh.LogIfErr(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		mq.SourceExhangeName,
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	mq.Errh.LogIfErr(err, "Failed to register a consumer")

	forever := make(chan bool)


	go func() {
		for d := range msgs {
			mq.Log.Info(d.ContentType)
			eb := &evented_core.EventBook{}
			err := proto.Unmarshal(d.Body, eb)
			mq.Errh.LogIfErr(err, "Failed to Unmarshal Event Book")
			uuid, err := evented_proto.ProtoToUUID(*eb.Cover.GetRoot())
			mq.Errh.LogIfErr(err, "Failed unparse UUID")
			mq.Log.Infof("Received a message: %s", uuid)
			response, err := mq.EventHandler.Handle(context.Background(), eb)
			if response != nil {
				err = mq.Sender.Handle(response)
				mq.Errh.LogIfErr(err, "Failed to send response eventbook to next transmission medium")
			}
		}
	}()

	mq.Log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}