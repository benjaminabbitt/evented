package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type Holder interface {
	AddEventBookChan(ebc chan *evented.EventBook)
	AddProjectorClient(pc projector.SyncProjectorTransporter)
	AddSagaTransporter(st saga.SyncSagaTransporter)
	GetTransports() []chan *evented.EventBook
	GetProjectors() []projector.SyncProjectorTransporter
	GetSaga() []saga.SyncSagaTransporter
}
