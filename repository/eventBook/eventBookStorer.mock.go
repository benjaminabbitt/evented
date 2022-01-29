package eventBook

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockEventBookRepository struct {
	mock.Mock
}

func (o *MockEventBookRepository) Get(ctx context.Context, id uuid.UUID) (book *core.EventBook, err error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*core.EventBook), args.Error(1)
}

func (o *MockEventBookRepository) Put(ctx context.Context, book *core.EventBook) error {
	args := o.Called(ctx, book)
	return args.Error(0)
}

func (o *MockEventBookRepository) GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *core.EventBook, err error) {
	args := o.Called(ctx, id, from, to)
	return args.Get(0).(*core.EventBook), args.Error(1)
}

func (o *MockEventBookRepository) GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *core.EventBook, err error) {
	args := o.Called(ctx, id, from)
	return args.Get(0).(*core.EventBook), args.Error(1)
}
