package memoryRepository

type MemoryRepository struct{
	storage map[string][]Event
}

func(repos MemoryRepository) Add(evt Event){
	repos.storage[evt.id] = append(repos.storage[evt.id],  evt)
}

func(repos MemoryRepository) Get(id string)(events []Event){
	return repos.storage[id]
}

func(repos MemoryRepository) GetTo(id string, to uint32)(events []Event){
	return repos.storage[id]
}

func(repos MemoryRepository) GetFromTo(id string, from uint32, to uint32)(events []Event){

}
