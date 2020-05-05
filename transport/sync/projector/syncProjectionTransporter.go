package projector

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error)
}
