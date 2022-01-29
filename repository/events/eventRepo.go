package events

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
)

type EventStorer interface {
	Add(ctx context.Context, id uuid.UUID, evt []*evented.EventPage) (err error)

	Get(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID) (err error)

	//Note that GetTo treats the to parameter as an exclusive upper bound, like Go slices.
	GetTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, to uint32) (err error)

	GetFrom(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32) (err error)

	//Note that GetTo treats the to parameter as an exclusive upper bound, like Go slices.
	GetFromTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32, to uint32) (err error)
}
