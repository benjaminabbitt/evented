package client

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
)

type BusinessClient interface {
	Handle(ctx context.Context, command *eventedcore.ContextualCommand) (events *eventedcore.EventBook, err error)
}
