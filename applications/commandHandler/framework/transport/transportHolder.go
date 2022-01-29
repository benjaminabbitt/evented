package transport

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
)

type TransportHolder interface {
	Add(i interface{}) error
	GetTransports() []chan *evented.EventBook
	GetProjectors() []evented.ProjectorClient
	GetSaga() []saga.SyncSagaTransporter
}
