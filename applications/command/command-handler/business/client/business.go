package client

import (
	"context"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
)

type BusinessClient interface {
	Handle(ctx context.Context, command *evented.ContextualCommand) (events *evented.EventBook, err error)
}
