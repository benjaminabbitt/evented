package receiver

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPDecodedMessage struct {
	Book *evented_core.EventBook
	Ack  func() error
	Nack func() error
}

type AMQPReceiver struct {
	SourceURL         string
	SourceExhangeName string
	SourceQueueName   string
	Log               *zap.SugaredLogger
	ch                *amqp.Channel
	queue             *amqp.Queue
	OutputChannel     chan<- AMQPDecodedMessage
	deliveryChan      <-chan amqp.Delivery
	conn              *amqp.Connection
}

func (o *AMQPReceiver) ListenForever() {
	forever := make(chan bool)

	go o.Listen()

	o.Log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (o *AMQPReceiver) Listen() {
	for delivery := range o.deliveryChan {
		eb, ack, nack := o.ExtractMessage(delivery)

		o.OutputChannel <- AMQPDecodedMessage{
			Book: eb,
			Ack:  ack,
			Nack: nack,
		}
	}
}

func (o *AMQPReceiver) ExtractMessage(delivery amqp.Delivery) (book *evented_core.EventBook, ack func() error, nack func() error) {
	book = &evented_core.EventBook{}
	o.Log.Info(delivery.ContentType)
	err := proto.Unmarshal(delivery.Body, book)
	if err != nil {
		o.Log.Error(err)
	}
	if book.Cover == nil {
		o.Log.Errorw("Cover is nil, this should not be possible here", "book", book, "cover", book.Cover)
	}
	uuid, err := evented_proto.ProtoToUUID(book.Cover.GetRoot())
	if err != nil {
		o.Log.Error(err)
	}
	o.Log.Infof("Received a message: %s", uuid)
	return book, func() error { return o.ch.Ack(delivery.DeliveryTag, false) }, func() error { return o.ch.Nack(delivery.DeliveryTag, false, true) }
}

func (o *AMQPReceiver) Connect() error {
	err := backoff.Retry(func() error {
		var err error
		var conn *amqp.Connection
		conn, err = amqp.Dial(o.SourceURL)
		o.conn = conn
		return err
	}, backoff.NewExponentialBackOff())

	if err != nil {
		o.Log.Error(err)
	}

	ch, err := o.conn.Channel()
	if err != nil {
		o.Log.Error(err)
		return err
	}
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
			true,  // durable
			false, // delete when unused
			false, // exclusive
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
