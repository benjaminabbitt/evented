package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
)

type SyncProjectorTransporter interface {
	HandleSync(ctx context.Context, in *evented.EventBook, opts ...grpc.CallOption) (*evented.Projection, error)
}
