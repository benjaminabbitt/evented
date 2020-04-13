package transport

import (
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{})
	GetTransports() []async.EventTransporter
	GetProjections() []projector.SyncProjectionTransporter
	GetSaga() []saga.SyncSagaTransporter
}
