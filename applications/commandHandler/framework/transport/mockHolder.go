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

func (o MockHolder) GetTransports() []async.Transport {
	args := o.Called()
	return args.Get(0).([]async.Transport)
}

func (o MockHolder) GetProjections() []projector.SyncProjection {
	args := o.Called()
	return args.Get(0).([]projector.SyncProjection)
}

func (o MockHolder) GetSaga() []saga.SyncSaga {
	args := o.Called()
	return args.Get(0).([]saga.SyncSaga)
}