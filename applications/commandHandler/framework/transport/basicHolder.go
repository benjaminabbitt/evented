package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"go.uber.org/zap"
	"reflect"
)

type BasicHolder struct {
	Log         *zap.SugaredLogger
	transports  []chan *evented_core.EventBook
	projections []projector.SyncProjectionTransporter
	sagas       []saga.SyncSagaTransporter
}

func (th BasicHolder) Add(i interface{}) {
	switch i.(type) {
	case chan *evented_core.EventBook:
		th.transports = append(th.transports, i.(chan *evented_core.EventBook))
	default:
		th.Log.Infow("Attempted to add non-transport type to transport BasicHolder.  This may be a synchronous-only transport, and may be OK.")
	}

	switch i.(type) {
	case projector.SyncProjectionTransporter:
		th.projections = append(th.projections, i.(projector.SyncProjectionTransporter))
	case saga.SyncSagaTransporter:
		th.sagas = append(th.sagas, i.(saga.SyncSagaTransporter))
	default:
		th.Log.Infow("Attempted to add non-synchronous type to transport BasicHolder.", "type", reflect.TypeOf(i).Name())
	}
}

func (th BasicHolder) GetTransports() []chan *evented_core.EventBook {
	return th.transports
}

func (th BasicHolder) GetProjections() []projector.SyncProjectionTransporter {
	return th.projections
}

func (th BasicHolder) GetSaga() []saga.SyncSagaTransporter {
	return th.sagas
}

func NewTransportHolder(log *zap.SugaredLogger) *BasicHolder {
	return &BasicHolder{Log: log}
}
