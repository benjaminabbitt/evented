package memory

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	evented_memory_ops "github.com/benjaminabbitt/evented/repository/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewEventRepoMemory(log *zap.SugaredLogger) (EventRepoMemory, error) {
	return EventRepoMemory{
		store: make(map[uuid.UUID][]*evented_core.EventPage),
		log:   log,
	}, nil
}

type EventRepoMemory struct {
	store map[uuid.UUID][]*evented_core.EventPage
	log   *zap.SugaredLogger
}

func (o EventRepoMemory) Add(ctx context.Context, id uuid.UUID, evt []*evented_core.EventPage) (err error) {
	//TODO, assertions
	storable, forced, remainder := evented_memory_ops.ExtractUntilFirstForced(evt)
	o.store[id] = append(o.store[id], storable...)
	storable = o.store[id]
	next := o.getNextSequence(id)

	if forced != nil {
		evented_memory_ops.SetSequence(forced, next)
		next += 1
		storable = append(storable, forced)
	}

	for _, ea := range remainder {
		evented_memory_ops.SetSequence(ea, next)
		next += 1
		storable = append(storable, ea)
	}

	o.store[id] = storable
	return nil
}

func (o EventRepoMemory) getLastSequence(id uuid.UUID) uint32 {
	row := o.store[id]
	var last uint32 = 0
	if len(row) >= 1 {
		last = row[len(row)-1].GetNum()
	}
	return last
}

func (o EventRepoMemory) getNextSequence(id uuid.UUID) uint32 {
	row := o.store[id]
	if len(row) == 0 {
		return 0
	} else {
		return o.getLastSequence(id) + 1
	}
}

func (o EventRepoMemory) Get(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID) (err error) {
	go o.send(evtChan, o.store[id])
	return nil
}

func (o EventRepoMemory) GetTo(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, to uint32) (err error) {
	return o.GetFromTo(ctx, evtChan, id, 0, to)
}
func (o EventRepoMemory) GetFrom(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, from uint32) (err error) {
	return o.GetFromTo(ctx, evtChan, id, from, o.getNextSequence(id))
}
func (o EventRepoMemory) GetFromTo(ctx context.Context, evtChan chan *evented_core.EventPage, id uuid.UUID, from uint32, to uint32) (err error) {
	slice := o.filterStore(id, from, to)
	go o.send(evtChan, slice)
	return nil
}

func (o EventRepoMemory) filterStore(id uuid.UUID, from uint32, to uint32) []*evented_core.EventPage {
	return o.store[id][from:to]
}

func (o EventRepoMemory) send(echan chan *evented_core.EventPage, events []*evented_core.EventPage) {
	for _, evt := range events {
		echan <- evt
	}
	close(echan)
}
