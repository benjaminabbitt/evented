package events

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
)

type EventStorer interface {
	//Add adds the event page to the data store
	Add(ctx context.Context, id uuid.UUID, evt []*evented.EventPage) (err error)

	//Get retrieves the event pages associated with id and loads them into the event channel
	Get(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID) (err error)

	//GetTo retrieves the event pages associated with the ID up to (and exclusive of) the upper bound and loads them into the event channel
	GetTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, to uint32) (err error)

	//GetFrom retrieves the event pages associated with the ID from (and inclusive of) the lower bound and loads them into the event channel
	GetFrom(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32) (err error)

	//GetFromTo retrieves the event pages associated with the ID from the (and inclusive of) the lower bound and to (and exclusive of) the upper bound.  It loads them into the event channel
	GetFromTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32, to uint32) (err error)
}
