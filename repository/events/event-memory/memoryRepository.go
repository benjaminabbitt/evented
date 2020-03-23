package event_memory

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

type MemoryRepository struct {
	storage map[string][]*evented_core.EventPage
}

func (repos MemoryRepository) Add(id uuid.UUID, ent []*evented_core.EventPage) (err error) {
	for _, event := range ent {
		var combinedEntity = repos.storage[id.String()]
		combinedEntity = append(combinedEntity, event)
		repos.storage[id.String()] = combinedEntity
	}
	return nil
}

func (repos MemoryRepository) Get(id uuid.UUID) (evts []*evented_core.EventPage, err error) {
	return repos.storage[id.String()], nil
}

func (repos MemoryRepository) GetTo(id uuid.UUID, to uint32) (evts []*evented_core.EventPage, err error) {
	unfiltered, err := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x *evented_core.EventPage) bool {
		return x.Sequence <= to
	}).([]*evented_core.EventPage)
	return filtered, nil
}

func (repos MemoryRepository) GetFrom(id uuid.UUID, from uint32) (evts []*evented_core.EventPage, err error) {
	unfiltered, _ := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x *evented_core.EventPage) bool {
		return x.Sequence >= from
	}).([]*evented_core.EventPage)
	return filtered, nil
}

func (repos MemoryRepository) GetFromTo(id uuid.UUID, from uint32, to uint32) (evts []*evented_core.EventPage, err error) {
	unfiltered, err := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x *evented_core.EventPage) bool {
		return x.Sequence >= from && x.Sequence <= to
	}).([]*evented_core.EventPage)
	return filtered, nil
}

func NewMemoryRepository() (repos MemoryRepository) {
	repos = MemoryRepository{}
	repos.storage = make(map[string][]*evented_core.EventPage)
	return repos
}
