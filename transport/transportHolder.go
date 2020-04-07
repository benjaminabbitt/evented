package transport

import (
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"reflect"
)

type Holder interface {
	Add(i interface{})
	GetTransports() []async.Transport
	GetProjections() []projector.SyncProjection
	GetSaga() []saga.SyncSaga
}

type BasicHolder struct {
	Log         *zap.SugaredLogger
	transports  []async.Transport
	projections []projector.SyncProjection
	sagas       []saga.SyncSaga
}

func (th *BasicHolder) Add(i interface{}) {
	switch i.(type) {
	case async.Transport:
		th.transports = append(th.transports, i.(async.Transport))
	default:
		th.Log.Infow("Attempted to add non-transport type to transport BasicHolder.  This may be a synchronous-only transport, and may be OK.")
	}

	switch i.(type) {
	case projector.SyncProjection:
		th.projections = append(th.projections, i.(projector.SyncProjection))
	case saga.SyncSaga:
		th.sagas = append(th.sagas, i.(saga.SyncSaga))
	default:
		th.Log.Infow("Attempted to add non-synchronous type to transport BasicHolder.", "type", reflect.TypeOf(i).Name())
	}
}

func (th *BasicHolder) GetTransports() []async.Transport {
	return th.transports
}

func (th *BasicHolder) GetProjections() []projector.SyncProjection {
	return th.projections
}

func (th *BasicHolder) GetSaga() []saga.SyncSaga {
	return th.sagas
}

func NewTransportHolder(log *zap.SugaredLogger) *BasicHolder {
	return &BasicHolder{Log: log}
}

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
