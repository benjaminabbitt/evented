package projector

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type SyncProjectionTransporter interface {
	HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error)
}
