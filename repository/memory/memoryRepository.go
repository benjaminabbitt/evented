package memoryRepository

import (
	"github.com/benjaminabbitt/evented/framework"
	"github.com/thoas/go-funk"
)
import "github.com/benjaminabbitt/evented/repository"

type MemoryRepository struct {
	storage map[string]repository.Entity
}

func (repos MemoryRepository) Add(ent framework.Event) (err error) {
	var combinedEntity = repos.storage[ent.Id]
	combinedEntity.Events = append(combinedEntity.Events, repository.ConvertFrameworkEventToStorageEvent(ent))
	repos.storage[ent.Id] = combinedEntity
	return nil
}

func (repos MemoryRepository) Get(id string) (events []framework.Event, err error) {
	events = funk.Map(repos.storage[id].Events, func(storageEvent repository.Event) framework.Event {
		return framework.Event{
			Id:       id,
			Sequence: storageEvent.Sequence,
			Details:  storageEvent.Details,
		}
	}).([]framework.Event)
	return events, nil
}

func (repos MemoryRepository) GetTo(id string, to uint32) (events []framework.Event, err error) {
	unfiltered, err := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x framework.Event) bool {
		return x.Sequence <= to
	}).([]framework.Event)
	return filtered, nil
}

func (repos MemoryRepository) GetFromTo(id string, from uint32, to uint32) (events []framework.Event, err error) {
	unfiltered, err := repos.Get(id)
	var filtered = funk.Filter(unfiltered, func(x framework.Event) bool {
		return x.Sequence >= from && x.Sequence <= to
	}).([]framework.Event)
	return filtered, nil
}

func NewMemoryRepository() (repos MemoryRepository) {
	repos = MemoryRepository{}
	repos.storage = make(map[string]repository.Entity)
	return repos
}
