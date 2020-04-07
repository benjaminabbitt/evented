package client

import evented_core "github.com/benjaminabbitt/evented/proto/core"

type BusinessClient interface {
	Handle(command *evented_core.ContextualCommand) (events *evented_core.EventBook, err error)
}