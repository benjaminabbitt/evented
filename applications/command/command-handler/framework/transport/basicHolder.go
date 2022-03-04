package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"go.uber.org/zap"
)

type BasicHolder struct {
	Log        *zap.SugaredLogger
	transports []chan *evented.EventBook
	projectors []projector.SyncProjectorTransporter
	sagas      []saga.SyncSagaTransporter
}

func (th *BasicHolder) AddEventBookChan(ebc chan *evented.EventBook) {
	th.transports = append(th.transports, ebc)
}

func (th *BasicHolder) AddProjectorClient(pc projector.SyncProjectorTransporter) {
	th.projectors = append(th.projectors, pc)
}

func (th *BasicHolder) AddSagaTransporter(st saga.SyncSagaTransporter) {
	th.sagas = append(th.sagas, st)
}

func (th *BasicHolder) GetTransports() []chan *evented.EventBook {
	return th.transports
}

func (th *BasicHolder) GetProjectors() []projector.SyncProjectorTransporter {
	return th.projectors
}

func (th *BasicHolder) GetSaga() []saga.SyncSagaTransporter {
	return th.sagas
}

func NewTransportHolder(log *zap.SugaredLogger) *BasicHolder {
	return &BasicHolder{Log: log}
}
