package events

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type EventRepository interface {
	Add(ctx context.Context, id uuid.UUID, evt []*evented_core.EventPage) (err error)
	Get(ctx context.Context, id uuid.UUID) (evt []*evented_core.EventPage, err error)
	GetTo(ctx context.Context, id uuid.UUID, to uint32) (evt []*evented_core.EventPage, err error)
	GetFrom(ctx context.Context, id uuid.UUID, from uint32) (evt []*evented_core.EventPage, err error)
	GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (evt []*evented_core.EventPage, err error)
}
