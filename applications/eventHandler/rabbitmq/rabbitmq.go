package rabbitmq

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQReceiver struct {
	SourceURL         string
	SourceExhangeName string
	SourceQueueName   string
	DestinationSink   map[string]evented_core.CommandHandlerClient
	Log               *zap.SugaredLogger
	EventHandler      evented_eventHandler.EventHandlerClient
}

func (mq *RabbitMQReceiver) Listen() {
	viper.GetString("")

	conn, err := amqp.Dial(mq.SourceURL)
	if err != nil {
		mq.Log.Error(err)
	}
	if conn != nil {
		defer conn.Close()
	}

	ch, err := conn.Channel()

	if ch != nil {
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

		q, err := ch.QueueDeclare(
			mq.SourceQueueName,
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			mq.Log.Error(err)
		}

		err = ch.QueueBind(
			q.Name,
			"",
			mq.SourceExhangeName,
			false,
			nil,
		)
		if err != nil {
			mq.Log.Error(err)
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				mq.Log.Info(d.ContentType)
				eb := &evented_core.EventBook{}
				err := proto.Unmarshal(d.Body, eb)
				if err != nil {
					mq.Log.Error(err)
				}
				uuid, err := evented_proto.ProtoToUUID(eb.Cover.GetRoot())
				if err != nil {
					mq.Log.Error(err)
				}
				mq.Log.Infof("Received a message: %s", uuid)
				response, err := mq.EventHandler.Handle(context.Background(), eb)
				if response != nil {
					chResponse, err := mq.DestinationSink[response.Cover.Domain].Record(context.Background(), response)
					if err != nil {
						mq.Log.Error(err)
					}
					mq.Log.Info(chResponse)
				}
			}
		}()

		mq.Log.Info(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}
}
