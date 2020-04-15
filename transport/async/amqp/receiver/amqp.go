package receiver

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPReceiver struct {
	SourceURL         string
	SourceExhangeName string
	SourceQueueName   string
	DestinationSink   map[string]evented_core.CommandHandlerClient
	Log               *zap.SugaredLogger
	EventHandler      evented_eventHandler.EventHandlerClient
	ch                *amqp.Channel
	queue             *amqp.Queue
	deliveryChan      <-chan amqp.Delivery
}

func (o *AMQPReceiver) ListenForever() {
	forever := make(chan bool)

	go func() {
		for delivery := range o.deliveryChan {
			eb := o.ExtractMessage(delivery)
			err := o.ProcessMessage(context.Background(), eb)
			if err != nil {
				o.Log.Error(err)
			}
		}
	}()

	o.Log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (o *AMQPReceiver) ProcessMessage(ctx context.Context, book *evented_core.EventBook) error {
	response, err := o.EventHandler.Handle(context.Background(), book)
	if err != nil {
		o.Log.Error(err)
	}
	if response != nil {
		chResponse, err := o.DestinationSink[response.Cover.Domain].Record(context.Background(), response)
		if err != nil {
			o.Log.Error(err)
		}
		o.Log.Info(chResponse)
	}
	return nil
}

func (o *AMQPReceiver) ExtractMessage(delivery amqp.Delivery) *evented_core.EventBook {
	o.Log.Info(delivery.ContentType)
	eb := &evented_core.EventBook{}
	err := proto.Unmarshal(delivery.Body, eb)
	if err != nil {
		o.Log.Error(err)
	}
	uuid, err := evented_proto.ProtoToUUID(eb.Cover.GetRoot())
	if err != nil {
		o.Log.Error(err)
	}
	o.Log.Infof("Received a message: %s", uuid)
	return eb
}

func (o *AMQPReceiver) Connect() error {
	conn, err := amqp.Dial(o.SourceURL)
	if err != nil {
		o.Log.Error(err)
		return err
	}

	ch, err := conn.Channel()
	o.ch = ch

	if ch != nil {
		err = ch.ExchangeDeclare(
			o.SourceExhangeName,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)

		q, err := ch.QueueDeclare(
			o.SourceQueueName,
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			o.Log.Error(err)
			return err
		}
		o.queue = &q

		err = ch.QueueBind(
			o.queue.Name,
			"",
			o.SourceExhangeName,
			false,
			nil,
		)
		if err != nil {
			o.Log.Error(err)
			return err
		}

		delivery, err := o.ch.Consume(
			o.queue.Name, // queue
			"",           // consumer
			true,         // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			o.Log.Error(err)
			return err
		}
		o.deliveryChan = delivery
	}
	return nil
}

func (o *AMQPReceiver) GetMessage(ctx context.Context) *evented_core.EventBook {
	var delivery amqp.Delivery
	delivery = <-o.deliveryChan
	msg := o.ExtractMessage(delivery)
	return msg
}
