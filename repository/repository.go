package repository


type EventRepository interface {
	Add(evt StorageEvent)
	Get(id string) (events []StorageEvent)
	GetTo(id string, to uint32) (events []StorageEvent)
	GetFromTo(id string, from uint32, to uint32) (events []StorageEvent)
}