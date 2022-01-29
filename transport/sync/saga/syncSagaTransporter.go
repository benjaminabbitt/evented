package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
)

type SyncSagaTransporter interface {
	HandleSync(ctx context.Context, evts *core.EventBook) (response *core.SynchronousProcessingResponse, err error)
}
