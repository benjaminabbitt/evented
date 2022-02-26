package sender

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type AmqpSuite struct {
	suite.Suite
	client *AMQPSender
	dait   *dockerTestSuite.DockerAssistedIntegrationTest
	log    *zap.SugaredLogger
}

func (o *AmqpSuite) SetupSuite() {
	o.log = support.Log()

	o.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := o.dait.CreateNewContainer("rabbitmq:3.9.13-alpine", []uint16{4369, 5671, 5672, 25672})
	if err != nil {
		o.log.Error(err)
	}
	time.Sleep(30 * time.Second)

	port, err := o.dait.GetPortMapping(5672)
	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)
	senderCh := make(chan evented.EventBook)
	o.client = NewAMQPSender(senderCh, url, "testExchange", o.log)
	err = o.client.Connect()
	if err != nil {
		o.log.Error(err)
	}
}

func (o *AmqpSuite) TearDownSuite() {
	err := o.dait.StopContainer()
	if err != nil {
		o.log.Error(err)
	}
}

func (o AmqpSuite) TestNoExceptionThrown() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented.EventPage{framework.NewEmptyEventPage(0, false)}, nil)
	err := o.client.Handle(eb)
	o.Assert().Nil(err)
}

func TestServerSuite(t *testing.T) {
	if !testing.Short() {
		suite.Run(t, new(AmqpSuite))
	}
}
