package amqp

import (
	"context"
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPSender struct {
	log          *zap.SugaredLogger
	ch           *amqp.Channel
	marshaller   *jsonpb.Marshaler
	exchangeName string
}

func (client AMQPSender) Handle(ctx context.Context, evts *evented_core.EventBook) (err error) {
	body, err := proto.Marshal(evts)
	client.log.Infow("Publishing ", "eventBook", support.StringifyEventBook(evts), "exchange", client.exchangeName)
	err = client.ch.Publish(
		client.exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: fmt.Sprintf("application/protobuf;proto=%T", *evts),
			Body:        []byte(body),
		})
	return nil
}

func NewAMQPSender(url string, exchangeName string, log *zap.SugaredLogger) *AMQPSender {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error(err)
	}
	log.Info("Connected to AMQP Broker")
	ch, err := conn.Channel()
	log.Info("Channel Formed")
	if err != nil {
		log.Error(err)
	}
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	log.Info("Exchange Declared")
	if err != nil {
		log.Error(err)
	}
	client := &AMQPSender{
		log:          log,
		exchangeName: exchangeName,
		ch:           ch,
		marshaller:   &jsonpb.Marshaler{},
	}
	return client
}
