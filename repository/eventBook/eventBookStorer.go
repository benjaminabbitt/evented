package eventBook

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
)

type Storer interface {
	Get(ctx context.Context, id uuid.UUID) (book *core.EventBook, err error)
	Put(ctx context.Context, book *core.EventBook) error
	GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *core.EventBook, err error)
	GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *core.EventBook, err error)
}
