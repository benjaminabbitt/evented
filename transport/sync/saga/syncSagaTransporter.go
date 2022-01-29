package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
)

type SyncSagaTransporter interface {
	HandleSync(ctx context.Context, evts *evented.EventBook) (response *evented.SynchronousProcessingResponse, err error)
}
