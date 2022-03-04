package sender

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/benjaminabbitt/evented/support"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPSender struct {
	log          *zap.SugaredLogger
	amqpch       *amqp.Channel
	conn         *amqp.Connection
	ch           chan *evented.EventBook
	exchangeName string
	url          string
}

func (o AMQPSender) Handle(evts *evented.EventBook) (err error) {
	body, err := proto.Marshal(evts)
	o.log.Infow("Publishing ", "eventBook", support.StringifyEventBook(evts), "exchange", o.exchangeName)
	err = backoff.Retry(func() error {
		return o.amqpch.Publish(
			o.exchangeName,
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: fmt.Sprintf("application/protobuf;proto=%T", evts),
				Body:        body,
			})
	}, backoff.NewExponentialBackOff())
	return nil
}

func (o AMQPSender) Run() {
	go func(ch chan *evented.EventBook) {
		for eb := range o.ch {
			err := o.Handle(eb)
			if err != nil {
				o.log.Error(err)
			}
		}
	}(o.ch)
}

func NewAMQPSender(ch chan *evented.EventBook, url string, exchangeName string, log *zap.SugaredLogger) *AMQPSender {
	client := &AMQPSender{
		log:          log,
		exchangeName: exchangeName,
		url:          url,
		ch:           ch,
	}
	return client
}

func (o *AMQPSender) Connect() error {
	err := backoff.Retry(func() error {
		var err error
		o.conn, err = amqp.Dial(o.url)
		return err
	}, backoff.NewExponentialBackOff())

	if err != nil {
		o.log.Error(err)
	}
	o.log.Info("Connected to AMQP Broker")
	ch, err := o.conn.Channel()
	if err != nil {
		o.log.Error(err)
		return err
	}
	o.amqpch = ch
	o.log.Info("Channel Formed")
	err = ch.ExchangeDeclare(
		o.exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	o.log.Info("Exchange Declared")
	if err != nil {
		o.log.Error(err)
		return err
	}
	return nil
}
