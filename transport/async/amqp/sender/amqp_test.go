package sender

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
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
	err := o.dait.CreateNewContainer("rabbitmq:3.8.3-alpine", []uint16{4369, 5671, 5672, 25672})
	if err != nil {
		o.log.Error(err)
	}
	time.Sleep(30 * time.Second)

	port, err := o.dait.GetPortMapping(5672)
	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)
	senderCh := make(chan *evented_core.EventBook)
	o.client = NewAMQPSender(senderCh, url, "testExchange", o.log)
	o.client.Connect()
}

func (o *AmqpSuite) TearDownSuite() {
	o.dait.StopContainer()
}

func (o AmqpSuite) TestNoExceptionThrown() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEmptyEventPage(0, false)}, nil)
	err := o.client.Handle(eb)
	o.Assert().Nil(err)
}

func TestServerSuite(t *testing.T) {
	if !testing.Short() {
		suite.Run(t, new(AmqpSuite))
	}
}
