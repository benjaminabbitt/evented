package sender

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
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
	err := o.dait.CreateNewContainer("rabbitmq", []uint16{4369, 5671, 5672, 25672})
	if err != nil {
		o.log.Error(err)
	}
	//time.Sleep(30 * time.Second)

	port, err := o.dait.GetPortMapping(5672)
	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", port)
	o.client = NewAMQPSender(url, "testExchange", o.log)
}

func (o *AmqpSuite) TearDownSuite() {
	o.dait.StopContainer()
}

func (o AmqpSuite) TestNoExceptionThrown() {
	id, _ := uuid.NewRandom()
	eb := framework.NewEventBook(id, "test", []*evented_core.EventPage{framework.NewEventPage(0, false, any.Any{})}, nil)
	err := o.client.Handle(eb)
	o.Assert().Nil(err)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(AmqpSuite))
}
