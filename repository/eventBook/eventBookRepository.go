package eventBook

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (book *evented_core.EventBook, err error)
	Put(ctx context.Context, book *evented_core.EventBook) error
	GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *evented_core.EventBook, err error)
	GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *evented_core.EventBook, err error)
}
