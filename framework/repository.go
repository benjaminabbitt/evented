package framework

type EventRepository interface {
	Add(evt Event, err error)
	Get(id string) (ent Event, err error)
	GetTo(id string, to uint32) (ent Event, err error)
	GetFromTo(id string, from uint32, to uint32) (ent Event, err error)
}
