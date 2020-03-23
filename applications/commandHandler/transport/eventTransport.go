package transport

import evented_core "github.com/benjaminabbitt/evented/proto/core"

type Transport interface {
	Handle(evts *evented_core.EventBook)(err error)
}