package projector

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type SyncProjection interface {
	HandleSync(evts *evented_core.EventBook) (projection *evented_core.Projection, err error)
}
