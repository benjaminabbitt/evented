package memoryRepository

import "github.com/thoas/go-funk"
import "github.com/benjaminabbitt/evented/repository"

type MemoryRepository struct{
	storage map[string][]repository.StorageEvent
}

func(repos MemoryRepository) Add(evt repository.StorageEvent){
	if repos.storage[evt.Id] == nil {
		repos.storage[evt.Id] = make([]repository.StorageEvent, 0)
	}
	repos.storage[evt.Id] = append(repos.storage[evt.Id],  evt)
}

func(repos MemoryRepository) Get(id string)(events []repository.StorageEvent){
	return repos.storage[id]
}

func(repos MemoryRepository) GetTo(id string, to uint32)(events []repository.StorageEvent){
	return funk.Filter(repos.storage[id], func(x repository.StorageEvent) bool {
		return x.Sequence <= to
	}).([]repository.StorageEvent)
}

func(repos MemoryRepository) GetFromTo(id string, from uint32, to uint32)(events []repository.StorageEvent){
	return funk.Filter(repos.storage[id], func(x repository.StorageEvent) bool {
		return x.Sequence >= from && x.Sequence <= to
	}).([]repository.StorageEvent)
}

func NewMemoryRepository()(repos MemoryRepository){
	repos = MemoryRepository{}
	repos.storage = make(map[string][]repository.StorageEvent)
	return repos
}


