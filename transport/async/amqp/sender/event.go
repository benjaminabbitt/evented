package sender

import "github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

type EventSender interface {
	Handle(evts *evented.EventBook) (err error)
	Run()
}
