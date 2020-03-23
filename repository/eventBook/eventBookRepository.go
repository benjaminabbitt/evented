package eventBook

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository struct {
	errh *evented.ErrLogger
	log *zap.SugaredLogger
	EventRepo    events.EventRepository
	SnapshotRepo snapshots.SnapshotRepo
	Domain       string
}

func (repo *Repository) Get(id uuid.UUID)(book evented_core.EventBook, err error) {
	snapshot, err := repo.SnapshotRepo.Get(id)
	repo.errh.LogIfErr(err, fmt.Sprintf("Failed to get snapshot for id %s", id))
	var from uint32 = 0
	if snapshot != nil {
		from = snapshot.Sequence
	}
	pages, err := repo.EventRepo.GetFrom(id, from)
	repo.errh.LogIfErr(err, fmt.Sprintf("Failed getting from page %d on id %s", from, id))
	return repo.makeEventBook(id, pages ,snapshot), nil
}

func (repo *Repository) Put(book evented_core.EventBook) error {
	root, err := evented_proto.ProtoToUUID(*book.Cover.Root)
	repo.errh.LogIfErr(err, "")
	err = repo.EventRepo.Add(root, book.Pages)
	repo.errh.LogIfErr(err, "Failed adding pages to repo")
	err = repo.SnapshotRepo.Put(root, book.Snapshot)
	repo.errh.LogIfErr(err, "Failed adding snapshot to repo")
	return err
}

func (repo *Repository) GetFromTo(id uuid.UUID, from uint32, to uint32)(book evented_core.EventBook, err error){
	eventPages, err := repo.EventRepo.GetFromTo(id, from, to)
	repo.errh.LogIfErr(err, fmt.Sprintf("Failed getting pages %d to %d on id %s", from, to, id))
	return repo.makeEventBook(id, eventPages, nil), nil
}

func (repo *Repository) GetFrom(id uuid.UUID, from uint32)(book evented_core.EventBook, err error){
	eventPages, err := repo.EventRepo.GetFrom(id, from)
	repo.errh.LogIfErr(err, fmt.Sprintf("Failed getting from page %d on id %s", from, id))
	return repo.makeEventBook(id, eventPages, nil), nil
}

func (repo *Repository) makeEventBook(root uuid.UUID, pages []*evented_core.EventPage, snapshot *evented_core.Snapshot)(book evented_core.EventBook){
	rootBytes, err := root.MarshalBinary()
	repo.errh.LogIfErr(err, fmt.Sprintf("Failed making Event Book"))
	protoRoot := &evented_core.UUID{
		Value: rootBytes,
	}
	cover := &evented_core.Cover{
		Domain: repo.Domain,
		Root: protoRoot,
	}
	book = evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}
	return book
}

