package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, evts *evented.EventBook, opts ...grpc.CallOption) (projection *evented.Projection, err error)
}
