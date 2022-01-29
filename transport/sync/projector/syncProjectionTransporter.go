package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, evts *core.EventBook) (projection *core.Projection, err error)
}
