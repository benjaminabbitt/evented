package client

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
)

type BusinessClient interface {
	Handle(ctx context.Context, command *core.ContextualCommand) (events *core.EventBook, err error)
}
