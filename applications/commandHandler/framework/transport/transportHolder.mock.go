package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/stretchr/testify/mock"
)

type MockHolder struct {
	mock.Mock
}

func (o MockHolder) Add(i interface{}) error {
	args := o.Called(i)
	return args.Error(0)
}

func (o MockHolder) GetTransports() []chan *evented.EventBook {
	args := o.Called()
	return args.Get(0).([]chan *evented.EventBook)
}

func (o MockHolder) GetProjectors() []evented.SyncProjectorTransporter {
	args := o.Called()
	return args.Get(0).([]evented.SyncProjectorTransporter)
}

func (o MockHolder) GetSaga() []saga.SyncSagaTransporter {
	args := o.Called()
	return args.Get(0).([]saga.SyncSagaTransporter)
}
