package sender

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/benjaminabbitt/evented/transport/async/amqp/receiver"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type AmqpSuite struct {
	suite.Suite
	sender       *sender.AMQPSender
	receiver     *receiver.AMQPReceiver
	dait         *dockerTestSuite.DockerAssistedIntegrationTest
	log          *zap.SugaredLogger
	exchangeName string
	queueName    string
}

func (o *AmqpSuite) SetupSuite() {
	o.log = support.Log()
	o.exchangeName = "testExchange"
	o.queueName = "testQueue"

	o.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := o.dait.CreateNewContainer("rabbitmq:management", []uint16{4369, 5671, 5672, 25672, 15672})
	if err != nil {
		o.log.Error(err)
	}
	time.Sleep(30 * time.Second)
	port, err := o.dait.GetPortMapping(5672)
	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)
	o.sender = sender.NewAMQPSender(url, o.exchangeName, o.log)
	o.sender.Connect()
	o.receiver = &receiver.AMQPReceiver{
		SourceURL:         url,
		SourceExhangeName: o.exchangeName,
		SourceQueueName:   o.queueName,
		DestinationSink:   nil,
		Log:               o.log,
		EventHandler:      nil,
	}
	o.receiver.Connect()
}

func (o *AmqpSuite) TearDownSuite() {
	o.dait.StopContainer()
}

func (o AmqpSuite) TestNoExceptionThrown() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEventPage(0, false, any.Any{})}, nil)
	err := o.sender.Handle(eb)
	o.Assert().Nil(err)
}

func (o AmqpSuite) TestSendAndReceive() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEventPage(0, false, any.Any{})}, nil)
	_ = o.sender.Handle(eb)
	time.Sleep(1 * time.Second)
	message := o.receiver.GetMessage(context.Background())
	o.Assert().Equal(eb.Cover.Domain, message.Cover.Domain)
	o.Assert().Equal(eb.Cover.Root.String(), message.Cover.Root.String())
	o.Assert().Equal(eb.Pages[0].Sequence, message.Pages[0].Sequence)

}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(AmqpSuite))
}
