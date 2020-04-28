package events

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type EventStorer interface {
	Add(ctx context.Context, id uuid.UUID, evt []*evented_core.EventPage) (err error)
	Get(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID) (err error)
	GetTo(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, to uint32) (err error)
	GetFrom(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, from uint32) (err error)
	GetFromTo(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, from uint32, to uint32) (err error)
	EstablishIndices() (err error)
}
