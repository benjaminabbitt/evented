package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, evts *evented.EventBook, opts ...grpc.CallOption) (projection *evented.Projection, err error)
}
