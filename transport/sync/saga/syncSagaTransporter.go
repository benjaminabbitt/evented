package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"google.golang.org/grpc"
)

type SyncSagaTransporter interface {
	HandleSync(ctx context.Context, in *evented.EventBook, opts ...grpc.CallOption) (*evented.SynchronousProcessingResponse, error)
}
