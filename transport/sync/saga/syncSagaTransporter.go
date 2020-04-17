package saga

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type SyncSagaTransporter interface {
	HandleSync(ctx context.Context, evts *evented_core.EventBook) (responseEvents []*evented_core.EventBook, err error)
}
