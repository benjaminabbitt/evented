package repository

import (
	"github.com/benjaminabbitt/evented/framework"
	"github.com/golang/protobuf/ptypes/any"
)

type Entity struct {
	Events []Event
}

type Event struct {
	Sequence uint32
	Details  any.Any
}

func ConvertFrameworkEventToStorageEvent(f framework.Event) Event {
	return Event{
		f.Sequence,
		f.Details,
	}
}
