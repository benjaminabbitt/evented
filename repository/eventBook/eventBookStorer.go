package eventBook

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
)

type Storer interface {
	Get(ctx context.Context, id uuid.UUID) (book *evented.EventBook, err error)
	Put(ctx context.Context, book *evented.EventBook) error
	GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *evented.EventBook, err error)
	GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *evented.EventBook, err error)
}
