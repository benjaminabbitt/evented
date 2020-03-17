package events

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type EventRepository interface {
	Add(id string, evt []*evented_core.EventPage) (err error)
	Get(id string) (evt []*evented_core.EventPage, err error)
	GetTo(id string, to uint32) (evt []*evented_core.EventPage, err error)
	GetFrom(id string, from uint32) (evt []*evented_core.EventPage, err error)
	GetFromTo(id string, from uint32, to uint32) (evt []*evented_core.EventPage, err error)
}

