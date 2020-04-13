package eventBook

import (
	"context"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RepositoryBasic struct {
	log          *zap.SugaredLogger
	EventRepo    events.EventStorer
	SnapshotRepo snapshots.SnapshotStorer
	Domain       string
}

func (o RepositoryBasic) Get(ctx context.Context, id uuid.UUID) (book *eventedcore.EventBook, err error) {
	snapshot, err := o.SnapshotRepo.Get(ctx, id)
	if err != nil {
		o.log.Error(err)
	}
	var from uint32 = 0
	if snapshot != nil {
		from = snapshot.Sequence
	}
	pages, err := o.EventRepo.GetFrom(ctx, id, from)
	if err != nil {
		o.log.Error(err)
	}
	return o.makeEventBook(id, pages, snapshot), nil
}

func (o RepositoryBasic) Put(ctx context.Context, book *eventedcore.EventBook) error {
	root, err := eventedproto.ProtoToUUID(book.Cover.Root)
	if err != nil {
		o.log.Error(err)
	}
	err = o.EventRepo.Add(ctx, root, book.Pages)
	if err != nil {
		o.log.Error(err)
	}
	err = o.SnapshotRepo.Put(ctx, root, book.Snapshot)
	if err != nil {
		o.log.Error(err)
	}
	return err
}

func (o RepositoryBasic) GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *eventedcore.EventBook, err error) {
	eventPages, err := o.EventRepo.GetFromTo(ctx, id, from, to)
	if err != nil {
		o.log.Error(err)
	}
	return o.makeEventBook(id, eventPages, nil), nil
}

func (o RepositoryBasic) GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *eventedcore.EventBook, err error) {
	eventPages, err := o.EventRepo.GetFrom(ctx, id, from)
	if err != nil {
		o.log.Error(err)
	}
	return o.makeEventBook(id, eventPages, nil), nil
}

func (o RepositoryBasic) makeEventBook(root uuid.UUID, pages []*eventedcore.EventPage, snapshot *eventedcore.Snapshot) (book *eventedcore.EventBook) {
	rootBytes, err := root.MarshalBinary()
	if err != nil {
		o.log.Error(err)
	}
	protoRoot := &eventedcore.UUID{
		Value: rootBytes,
	}
	cover := &eventedcore.Cover{
		Domain: o.Domain,
		Root:   protoRoot,
	}
	book = &eventedcore.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}
	return book
}
