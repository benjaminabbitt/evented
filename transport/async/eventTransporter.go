package async

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
)

type EventTransporter interface {
	Handle(ctx context.Context, evts *evented_core.EventBook) (err error)
}
