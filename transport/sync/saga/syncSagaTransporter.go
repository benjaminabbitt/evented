package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
)

type SyncSagaTransporter interface {
	HandleSync(ctx context.Context, evts *evented_core.EventBook) (response *evented_core.SynchronousProcessingResponse, err error)
}
