package client

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type BusinessClient interface {
	Handle(ctx context.Context, command *evented_core.ContextualCommand) (events *evented_core.EventBook, err error)
}
