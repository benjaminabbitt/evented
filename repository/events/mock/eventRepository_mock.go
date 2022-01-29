package mock

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type EventRepository struct {
	mock.Mock
}

func (m EventRepository) Add(ctx context.Context, id uuid.UUID, evt []*evented.EventPage) (err error) {
	args := m.Called(ctx, id, evt)
	return args.Error(0)
}

func (m EventRepository) Get(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID) (err error) {
	args := m.Called(ctx, evtChan, id)
	return args.Error(0)
}

func (m EventRepository) GetTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, to uint32) (err error) {
	args := m.Called(ctx, evtChan, id, to)
	return args.Error(0)
}

func (m EventRepository) GetFrom(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32) (err error) {
	args := m.Called(ctx, evtChan, id, from)
	return args.Error(0)
}

func (m EventRepository) GetFromTo(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32, to uint32) (err error) {
	args := m.Called(ctx, evtChan, id, from, to)
	return args.Error(0)
}

func (m EventRepository) EstablishIndices() error {
	args := m.Called()
	return args.Error(0)
}
