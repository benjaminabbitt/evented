package amqp

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/golang/protobuf/jsonpb"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type AMQPClient struct {
	q amqp.Queue
	ch amqp.Channel
	marshaller jsonpb.Marshaler
}

func (c *AMQPClient) Project(evts *evented_core.EventBook) (err error){
	body, err := c.marshaller.MarshalToString(evts)
	failOnError(err, "Failed to serialize Event Book")
	err = c.ch.Publish(
		"",
		c.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/json",
			Body:            []byte(body),
		})
}


func NewAMQPClient(url string, queueName string) *AMQPClient {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to declare queue")
	return &AMQPClient{
		q:  q,
		ch: *ch,
		marshaller: jsonpb.Marshaler{},
	}
}

