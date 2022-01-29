package async

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
)

type EventTransporter interface {
	Handle(ctx context.Context, evts *core.EventBook) (err error)
}
