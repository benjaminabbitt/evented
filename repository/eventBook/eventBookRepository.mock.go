package eventBook

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockEventBookRepository struct {
	mock.Mock
}

func (o MockEventBookRepository) Get(id uuid.UUID) (book evented_core.EventBook, err error) {
	args := o.Called(id)
	return args.Get(0).(evented_core.EventBook), args.Error(1)
}

func (o MockEventBookRepository) Put(book evented_core.EventBook) error {
	args := o.Called(book)
	return args.Error(0)
}

func (o MockEventBookRepository) GetFromTo(id uuid.UUID, from uint32, to uint32) (book evented_core.EventBook, err error) {
	args := o.Called(id, from, to)
	return args.Get(0).(evented_core.EventBook), args.Error(1)
}

func (o MockEventBookRepository) GetFrom(id uuid.UUID, from uint32) (book evented_core.EventBook, err error) {
	args := o.Called(id, from)
	return args.Get(0).(evented_core.EventBook), args.Error(1)
}
