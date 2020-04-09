package framework

import core "github.com/benjaminabbitt/evented/proto/core"

type BusinessLogic interface {
	LoadSnapshot(snapshot core.Snapshot)
	LoadEvents(events core.EventBook)
	ProcessCommand(book core.CommandBook) core.EventBook
}