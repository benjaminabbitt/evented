package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	mock_evented "github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/mocks"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BasicHolderSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	holder *BasicHolder
}

func (suite *BasicHolderSuite) SetupTest() {
	suite.holder = NewTransportHolder(support.Log())
	suite.ctrl = gomock.NewController(suite.T())
}

func (suite *BasicHolderSuite) TestSyncProjectorHandling() {
	projectorClient := mock_evented.NewMockProjectorClient(suite.ctrl)
	projectorSet := []projector.SyncProjectorTransporter{projectorClient}
	suite.holder.AddProjectorClient(projectorClient)
	suite.Assert().Equal(projectorSet, suite.holder.GetProjectors())
}

func (suite *BasicHolderSuite) TestSyncSagaHandling() {
	sagaClient := mock_evented.NewMockSagaClient(suite.ctrl)
	sagaSet := []saga.SyncSagaTransporter{sagaClient}
	suite.holder.AddSagaTransporter(sagaClient)
	suite.Assert().Equal(sagaSet, suite.holder.GetSaga())
}

func (suite *BasicHolderSuite) TestTransportHandling() {
	ch := make(chan *evented.EventBook)
	chSet := []chan *evented.EventBook{ch}
	suite.holder.AddEventBookChan(ch)
	suite.Assert().Equal(chSet, suite.holder.GetTransports())
}

func TestBasicHolderSuite(t *testing.T) {
	suite.Run(t, new(BasicHolderSuite))
}
