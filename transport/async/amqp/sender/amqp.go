package sender

import (
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"time"
)

type AMQPSender struct {
	log          *zap.SugaredLogger
	amqpch       *amqp.Channel
	conn         *amqp.Connection
	ch           chan *evented_core.EventBook
	exchangeName string
	url          string
}

func (o AMQPSender) Handle(evts *evented_core.EventBook) (err error) {
	body, err := proto.Marshal(evts)
	o.log.Infow("Publishing ", "eventBook", support.StringifyEventBook(evts), "exchange", o.exchangeName)
	err = o.amqpch.Publish(
		o.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: fmt.Sprintf("application/protobuf;proto=%T", *evts),
			Body:        []byte(body),
		})
	return nil
}

func (o AMQPSender) Run() {
	go func(ch chan *evented_core.EventBook) {
		for eb := range o.ch {
			err := o.Handle(eb)
			if err != nil {
				o.log.Error(err)
			}
		}
	}(o.ch)
}

func NewAMQPSender(ch chan *evented_core.EventBook, url string, exchangeName string, log *zap.SugaredLogger) *AMQPSender {
	client := &AMQPSender{
		log:          log,
		exchangeName: exchangeName,
		url:          url,
		ch:           ch,
	}
	return client
}

///TODO: extrapolate to something more generic
func (o *AMQPSender) connectWithBackoff() error {
	var conn *amqp.Connection
	var err error
	var count uint8
	var max int = 1000
	var min int = 0
	for {
		conn, err = amqp.Dial(o.url)
		if err == nil {
			break
		}
		randOffset := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
		primaryTime := time.Duration(int(math.Pow(2, float64(count)))*1000) * time.Millisecond
		time.Sleep(primaryTime + randOffset)
		count++
	}
	o.conn = conn
	return nil
}

func (o *AMQPSender) Connect() error {
	o.connectWithBackoff()
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
