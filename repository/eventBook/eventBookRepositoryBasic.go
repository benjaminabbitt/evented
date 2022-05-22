package eventBook

import (
	"context"
	actx2 "github.com/benjaminabbitt/evented/applications/command/command-handler/actx"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func MakeRepositoryBasic(actx *actx2.BasicCommandHandlerApplicationContext, eventRepo events.EventStorer, snapshotRepo snapshots.SnapshotStorer) *RepositoryBasic {
	return &RepositoryBasic{
		log:          actx.Log(),
		EventRepo:    eventRepo,
		SnapshotRepo: snapshotRepo,
		Domain:       actx.Config.GetString(configuration.Domain),
	}
}

type RepositoryBasic struct {
	log          *zap.SugaredLogger
	EventRepo    events.EventStorer
	SnapshotRepo snapshots.SnapshotStorer
	Domain       string
}

func (o RepositoryBasic) Get(ctx context.Context, id uuid.UUID) (book *evented.EventBook, err error) {
	eventPageChannel := make(chan *evented.EventPage, 10)
	snapshot, err := o.SnapshotRepo.Get(ctx, id)
	if err != nil {
		o.log.Error(err)
	}
	var from uint32 = 0
	if snapshot != nil {
		from = snapshot.Sequence
	}
	err = o.EventRepo.GetFrom(ctx, eventPageChannel, id, from)
	if err != nil {
		o.log.Error(err)
	}
	var pages []*evented.EventPage
	for {
		page, more := <-eventPageChannel
		if !more {
			break
		}
		pages = append(pages, page)
	}
	return o.makeEventBook(id, pages, snapshot), nil
}

func (o RepositoryBasic) Put(ctx context.Context, book *evented.EventBook) error {
	root, err := eventedproto.ProtoToUUID(book.Cover.Root)
	if err != nil {
		o.log.Error(err)
	}
	err = o.EventRepo.Add(ctx, root, book.Pages)
	if err != nil {
		o.log.Error(err)
	}
	if book.Snapshot != nil {
		err = o.SnapshotRepo.Put(ctx, root, book.Snapshot)
		if err != nil {
			o.log.Error(err)
		}
	}
	return err
}

func (o RepositoryBasic) GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (book *evented.EventBook, err error) {
	ch := make(chan *evented.EventPage)
	err = o.EventRepo.GetFromTo(ctx, ch, id, from, to)
	if err != nil {
		o.log.Error(err)
	}
	var pages []*evented.EventPage
	for page := range ch {
		pages = append(pages, page)
	}
	return o.makeEventBook(id, pages, nil), nil
}

func (o RepositoryBasic) GetFrom(ctx context.Context, id uuid.UUID, from uint32) (book *evented.EventBook, err error) {
	ch := make(chan *evented.EventPage)
	err = o.EventRepo.GetFrom(ctx, ch, id, from)
	if err != nil {
		o.log.Error(err)
	}
	var pages []*evented.EventPage
	for page := range ch {
		pages = append(pages, page)
	}
	return o.makeEventBook(id, pages, nil), nil
}

func (o RepositoryBasic) makeEventBook(root uuid.UUID, pages []*evented.EventPage, snapshot *evented.Snapshot) (book *evented.EventBook) {
	rootBytes, err := root.MarshalBinary()
	if err != nil {
		o.log.Error(err)
	}
	protoRoot := &evented.UUID{
		Value: rootBytes,
	}
	cover := &evented.Cover{
		Domain: o.Domain,
		Root:   protoRoot,
	}
	book = &evented.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}
	return book
}
