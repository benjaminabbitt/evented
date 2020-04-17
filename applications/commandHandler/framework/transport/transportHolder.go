package transport

import (
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{})
	GetTransports() []chan *ContextualEventBook
	GetProjections() []projector.SyncProjectionTransporter
	GetSaga() []saga.SyncSagaTransporter
}
