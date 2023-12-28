package sender

import (
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
)

type EventSender interface {
	Handle(evts *evented.EventBook) (err error)
	Run()
}
