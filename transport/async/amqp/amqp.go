package amqp

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Client struct {
	errh         *evented.ErrLogger
	log          *zap.SugaredLogger
	ch           *amqp.Channel
	marshaller   *jsonpb.Marshaler
	exchangeName string
}

func (client Client) Handle(evts *evented_core.EventBook) (err error) {
	body, err := proto.Marshal(evts)
	client.errh.LogIfErr(err, "Failed to serialize Event Book")
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

func NewAMQPClient(url string, exchangeName string, log *zap.SugaredLogger, errh *evented.ErrLogger) *Client {
	conn, err := amqp.Dial(url)
	errh.LogIfErr(err, "Failed to connect")
	log.Info("Connected to AMQP Broker")
	ch, err := conn.Channel()
	log.Info("Channel Formed")
	errh.LogIfErr(err, "Failed to open channel")
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
	errh.LogIfErr(err, "Failed to declare exchange")
	client := &Client{
		log:          log,
		exchangeName: exchangeName,
		ch:           ch,
		marshaller:   &jsonpb.Marshaler{},
	}
	return client
}
