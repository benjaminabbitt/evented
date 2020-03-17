package event_memory

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/thoas/go-funk"
)

type MemoryRepository struct {
	storage map[string][]*evented_core.EventPage
}

func (repos MemoryRepository) Add(id string, ent []*evented_core.EventPage) (err error) {
	for _, event := range ent {
		var combinedEntity = repos.storage[id]
		combinedEntity = append(combinedEntity, event)
		repos.storage[id] = combinedEntity
	}
	return nil
}

func (repos MemoryRepository) Get(id string) (evts []*evented_core.EventPage, err error) {
	return repos.storage[id], nil
}

func (repos MemoryRepository) GetTo(id string, to uint32) (evts []*evented_core.EventPage, err error) {
	unfiltered, err := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x *evented_core.EventPage) bool {
		return x.Sequence <= to
	}).([]*evented_core.EventPage)
	return filtered, nil
}

func (repos MemoryRepository) GetFrom(id string, from uint32) (evts []*evented_core.EventPage, err error) {
	unfiltered, _ := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x *evented_core.EventPage) bool {
		return x.Sequence >= from
	}).([]*evented_core.EventPage)
	return filtered, nil
}

func (repos MemoryRepository) GetFromTo(id string, from uint32, to uint32) (evts []*evented_core.EventPage, err error) {
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
