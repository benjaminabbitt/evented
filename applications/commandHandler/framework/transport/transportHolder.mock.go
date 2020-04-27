package transport

import (
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/stretchr/testify/mock"
)

type MockHolder struct {
	mock.Mock
}

func (o MockHolder) Add(i interface{}) {
	o.Called(i)
}

func (o MockHolder) GetTransports() []async.EventTransporter {
	args := o.Called()
	return args.Get(0).([]async.EventTransporter)
}

func (o MockHolder) GetProjections() []projector.SyncProjectorTransporter {
	args := o.Called()
	return args.Get(0).([]projector.SyncProjectorTransporter)
}

func (o MockHolder) GetSaga() []saga.SyncSagaTransporter {
	args := o.Called()
	return args.Get(0).([]saga.SyncSagaTransporter)
}