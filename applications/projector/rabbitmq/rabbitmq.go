package rabbitmq

import (
	"github.com/benjaminabbitt/evented"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)


type RabbitMQReceiver struct {
	Errh *evented.ErrLogger
	Log  *zap.SugaredLogger
}

func (mq *RabbitMQReceiver) Listen(){
	exchangeName := "evented_evented"

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	mq.Errh.LogIfErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mq.Errh.LogIfErr(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	mq.Errh.LogIfErr(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"testQueue", // name
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
		exchangeName,
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
			eb := &evented_core.EventBook{}
			proto.Unmarshal(d.Body, eb)
			mq.Log.Infof("Received a message: %s", eb.Cover.GetRoot())

		}
	}()

	mq.Log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}