package amqp

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func (s *AmqpSuite) SetupSuite() {
	s.log = support.Log()

	s.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := s.dait.CreateNewContainer("rabbitmq", []uint16{4369, 5671, 5672, 25672})
	if err != nil {
		s.log.Error(err)
	}

	port, err := s.dait.GetPortMapping(5672)
	s.client = &AMQPSender{}
}

func (s AmqpSuite) TestExtract0() {
	id, _ := uuid.NewRandom()
	log.Info(id.String())
	log.Info(s.client.extractTopicalElement(0, id))
	log.Info(s.client.extractTopicalElement(1, id))
	log.Info(s.client.extractTopicalElement(2, id))
	log.Info(s.client.extractTopicalElement(3, id))
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(AmqpSuite))
}
