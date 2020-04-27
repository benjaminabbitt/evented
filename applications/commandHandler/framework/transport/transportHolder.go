package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{})
	GetTransports() []chan *evented_core.EventBook
	GetProjections() []projector.SyncProjectorTransporter
	GetSaga() []saga.SyncSagaTransporter
}
