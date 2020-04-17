package transport

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type ContextualEventBook struct {
	EventBook *evented_core.EventBook
	Context   context.Context
}
