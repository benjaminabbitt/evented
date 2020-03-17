package framework

import (
evented_core "github.com/benjaminabbitt/evented/proto/core"
"github.com/golang/protobuf/ptypes/any"
"time"
)

func NewEventPage(sequence uint32, sync bool, eventDetails any.Any) *evented_core.EventPage {
	return &evented_core.EventPage{
		Sequence:  sequence,
		Synchronous: sync,
		CreatedAt: time.Now().Format(time.RFC3339),
		Event:     &eventDetails,
	}
}

func NewEventBook(id string, domain string, events []*evented_core.EventPage, snapshot *evented_core.Snapshot) *evented_core.EventBook{
	return &evented_core.EventBook{
		Cover:    &evented_core.Cover{
			Domain: domain,
			Root:     id,
		},
		Pages:    events,
		Snapshot: snapshot,
	}
}
