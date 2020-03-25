package transport

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type SyncSaga interface{
	HandleSync(evts *evented_core.EventBook)(responseEvents *evented_core.EventBook, err error)
}
