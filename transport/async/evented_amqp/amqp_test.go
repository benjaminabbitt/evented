package evented_amqp

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AmqpSuite struct {
	suite.Suite
	client *Client
}

func (s *AmqpSuite) SetupSuite(){
	s.client = &Client{}
}

func (s AmqpSuite) TestExtract0(){
	id, _  := uuid.NewRandom()
	log.Info(id.String())
	log.Info(s.client.extractTopicalElement(0 , id))
	log.Info(s.client.extractTopicalElement(1 , id))
	log.Info(s.client.extractTopicalElement(2 , id))
	log.Info(s.client.extractTopicalElement(3 , id))
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(AmqpSuite))
}