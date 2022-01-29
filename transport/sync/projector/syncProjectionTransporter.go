package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, evts *evented.EventBook) (projection *evented.Projection, err error)
}
