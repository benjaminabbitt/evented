package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
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
	projectorClient := evented.MockProjectorClient{}
	projectorSet := []evented.SyncProjectorTransporter{projectorClient}
	err := o.holder.Add(projectorClient)
	o.Assert().Equal(projectorSet, o.holder.GetProjectors())
	o.Assert().NoError(err)
}

func (o *BasicHolderSuite) TestSyncSagaHandling() {
	sagaClient := &saga.MockSagaClient{}
	sagaSet := []saga.SyncSagaTransporter{sagaClient}
	err := o.holder.Add(sagaClient)
	o.Assert().Equal(sagaSet, o.holder.GetSaga())
	o.Assert().NoError(err)
}

func (o *BasicHolderSuite) TestTransportHandling() {
	ch := make(chan *evented.EventBook)
	chSet := []chan *evented.EventBook{ch}
	err := o.holder.Add(ch)
	o.Assert().Equal(chSet, o.holder.GetTransports())
	o.Assert().NoError(err)
}

type empty struct{}

func (o *BasicHolderSuite) TestInvalidType() {
	foo := empty{}
	err := o.holder.Add(foo)
	o.Assert().Error(err)
}

func TestBasicHolderSuite(t *testing.T) {
	suite.Run(t, new(BasicHolderSuite))
}
