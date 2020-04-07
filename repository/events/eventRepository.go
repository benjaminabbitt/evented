package events

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type EventRepository interface {
	Add(id uuid.UUID, evt []*evented_core.EventPage) (err error)
	Get(id uuid.UUID) (evt []*evented_core.EventPage, err error)
	GetTo(id uuid.UUID, to uint32) (evt []*evented_core.EventPage, err error)
	GetFrom(id uuid.UUID, from uint32) (evt []*evented_core.EventPage, err error)
	GetFromTo(id uuid.UUID, from uint32, to uint32) (evt []*evented_core.EventPage, err error)
}
