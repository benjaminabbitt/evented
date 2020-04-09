package async

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type Transport interface {
	Handle(ctx context.Context, evts *evented_core.EventBook) (err error)
}
