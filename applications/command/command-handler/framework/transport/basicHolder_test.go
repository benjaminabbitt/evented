package transport

import (
	"github.com/benjaminabbitt/evented/mocks"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BasicHolderSuite struct {
	suite.Suite
	holder *BasicHolder
}

func (o *BasicHolderSuite) SetupTest() {
	o.holder = NewTransportHolder(support.Log())
}

func (o *BasicHolderSuite) TestSyncProjectorHandling() {
	projectorClient := &mocks.ProjectorClient{}
	projectorSet := []projector.SyncProjectorTransporter{projectorClient}
	o.holder.AddProjectorClient(projectorClient)
	o.Assert().Equal(projectorSet, o.holder.GetProjectors())
}

func (o *BasicHolderSuite) TestSyncSagaHandling() {
	sagaClient := &mocks.SyncSagaTransporter{}
	sagaSet := []saga.SyncSagaTransporter{sagaClient}
	o.holder.AddSagaTransporter(sagaClient)
	o.Assert().Equal(sagaSet, o.holder.GetSaga())
}

func (o *BasicHolderSuite) TestTransportHandling() {
	ch := make(chan *evented.EventBook)
	chSet := []chan *evented.EventBook{ch}
	o.holder.AddEventBookChan(ch)
	o.Assert().Equal(chSet, o.holder.GetTransports())
}

func TestBasicHolderSuite(t *testing.T) {
	suite.Run(t, new(BasicHolderSuite))
}
