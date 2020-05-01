package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{}) error
	GetTransports() []chan *evented_core.EventBook
	GetProjectors() []projector.SyncProjectorTransporter
	GetSaga() []saga.SyncSagaTransporter
}
