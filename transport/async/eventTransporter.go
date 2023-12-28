package async

import (
	"context"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
)

type EventTransporter interface {
	Handle(ctx context.Context, evts *evented.EventBook) (err error)
}
