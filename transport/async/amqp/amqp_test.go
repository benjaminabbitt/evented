package sender

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type AmqpSuite struct {
	suite.Suite
	sender        *sender.AMQPSender
	receiver      *receiver.AMQPReceiver
	dait          *dockerTestSuite.DockerAssistedIntegrationTest
	log           *zap.SugaredLogger
	outputChannel chan receiver.AMQPDecodedMessage
	exchangeName  string
	queueName     string
}

func (o *AmqpSuite) SetupSuite() {
	o.log = support.Log()
	o.exchangeName = "testExchange"
	o.queueName = "testQueue"
	o.outputChannel = make(chan receiver.AMQPDecodedMessage, 10)

	o.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := o.dait.CreateNewContainer("rabbitmq:3.8.3-alpine", []uint16{4369, 5671, 5672, 25672, 15672})
	if err != nil {
		o.log.Error(err)
	}
	port, err := o.dait.GetPortMapping(5672)
	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)
	ch := make(chan *evented_core.EventBook, 10)
	o.sender = sender.NewAMQPSender(ch, url, o.exchangeName, o.log)
	err = o.sender.Connect()
	if err != nil {
		o.log.Error(err)
	}

	o.receiver = &receiver.AMQPReceiver{
		SourceURL:         url,
		SourceExhangeName: o.exchangeName,
		SourceQueueName:   o.queueName,
		Log:               o.log,
		OutputChannel:     o.outputChannel,
	}
	err = o.receiver.Connect()
	if err != nil {
		o.log.Error(err)
	}
}

func (o *AmqpSuite) TearDownSuite() {
	o.dait.StopContainer()
}

func (o AmqpSuite) TestNoExceptionThrown() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEmptyEventPage(0, false)}, nil)
	err := o.sender.Handle(eb)
	o.Assert().Nil(err)
}

func (o AmqpSuite) TestSendAndReceive() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEmptyEventPage(0, false)}, nil)
	go o.receiver.Listen()
	_ = o.sender.Handle(eb)
	message := <-o.outputChannel
	o.Assert().Equal(eb.Cover.Domain, message.Book.Cover.Domain)
	o.Assert().Equal(eb.Cover.Root.String(), message.Book.Cover.Root.String())
	o.Assert().Equal(eb.Pages[0].Sequence, message.Book.Pages[0].Sequence)

}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(AmqpSuite))
}
