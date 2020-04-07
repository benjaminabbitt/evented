package mock

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type EventRepository struct {
	mock.Mock
}

func (m EventRepository) Add(id uuid.UUID, evt []*evented_core.EventPage) (err error) {
	args := m.Called(id, evt)
	return args.Error(0)
}

func (m EventRepository) Get(id uuid.UUID) (evt []*evented_core.EventPage, err error) {
	args := m.Called(id)
	return args.Get(0).([]*evented_core.EventPage), args.Error(1)
}

func (m EventRepository) GetTo(id uuid.UUID, to uint32) (evt []*evented_core.EventPage, err error) {
	args := m.Called(id, to)
	return args.Get(0).([]*evented_core.EventPage), args.Error(1)
}

func (m EventRepository) GetFrom(id uuid.UUID, from uint32) (evt []*evented_core.EventPage, err error) {
	args := m.Called(id, from)
	return args.Get(0).([]*evented_core.EventPage), args.Error(1)
}

func (m EventRepository) GetFromTo(id uuid.UUID, from uint32, to uint32) (evt []*evented_core.EventPage, err error) {
	args := m.Called(id, from, to)
	return args.Get(0).([]*evented_core.EventPage), args.Error(1)
}
