package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type Saga interface {
	SendSync(evts *evented_core.EventBook)(responseEvents *evented_core.EventBook, err error)
	Send(evts *evented_core.EventBook)(err error)
}
