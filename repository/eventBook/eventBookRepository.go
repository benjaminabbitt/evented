package eventBook

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/snapshots"
)

type Repository struct {
	EventRepo    events.EventRepository
	SnapshotRepo snapshots.SnapshotRepo
	Domain       string
}

func (repo *Repository) Get(id string)(book evented_core.EventBook, err error) {
	snapshot, err := repo.SnapshotRepo.Get(id)
	var pages []*evented_core.EventPage
	if err == nil {
		pages, err = repo.EventRepo.GetFrom(id, snapshot.Sequence)
	}
	cover := evented_core.Cover{
		Domain: repo.Domain,
		Root:   id,
	}
	book =  evented_core.EventBook{
		Cover:    &cover,
		Pages:    pages,
		Snapshot: snapshot,
	}
	return book, nil
}

func (repo *Repository) Put(book evented_core.EventBook) error {
	err := repo.EventRepo.Add(book.Cover.Root, book.Pages)
	if err == nil {
		err = repo.SnapshotRepo.Put(book.Cover.Root, book.Snapshot)
	}
	return err
}
