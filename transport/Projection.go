package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type Projection interface{
	ProjectSync(evts *evented_core.EventBook)(projection *evented_core.Projection, err error)
	Project(evts *evented_core.EventBook) (err error)
}