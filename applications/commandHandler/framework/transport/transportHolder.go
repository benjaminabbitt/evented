package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{}) error
	GetTransports() []chan *core.EventBook
	GetProjectors() []projector.SyncProjectorTransporter
	GetSaga() []saga.SyncSagaTransporter
}
