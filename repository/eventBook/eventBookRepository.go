package eventBook

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type Repository interface {
	Get(id uuid.UUID) (book evented_core.EventBook, err error)
	Put(book evented_core.EventBook) error
	GetFromTo(id uuid.UUID, from uint32, to uint32) (book evented_core.EventBook, err error)
	GetFrom(id uuid.UUID, from uint32) (book evented_core.EventBook, err error)
}

