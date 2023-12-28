package async

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
)

type EventTransporter interface {
	Handle(ctx context.Context, evts *evented.EventBook) (err error)
}
