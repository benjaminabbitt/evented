package eventBook

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	erMock "github.com/benjaminabbitt/evented/repository/events/mock"
	ssMock "github.com/benjaminabbitt/evented/repository/snapshots/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventBookRepositorySuite struct {
	suite.Suite
	domain              string
	id                  uuid.UUID
	eventBookRepository RepositoryBasic
	snapshotRepository  ssMock.SnapshotRepo
	eventRepository     erMock.EventRepository
	book                *evented_core.EventBook
}

func (o *EventBookRepositorySuite) SetupTest() {
	id, _ := uuid.NewRandom()
	o.id = id
	o.domain = "test"

	o.eventRepository = erMock.EventRepository{}
	o.snapshotRepository = ssMock.SnapshotRepo{}

	o.eventBookRepository = RepositoryBasic{
		EventRepo:    &o.eventRepository,
		SnapshotRepo: &o.snapshotRepository,
		Domain:       "test",
	}
}

//func (o *EventBookRepositorySuite) Test_Put() {
//	id, _ := uuid.NewRandom()
//	pid := evented_proto.UUIDToProto(id)
//
//	cover := &evented_core.Cover{
//		Domain: "testPut",
//		Root:   &pid,
//	}
//
//	pages := []*evented_core.EventPage{
//		&evented_core.EventPage{
//			Sequence: &evented_core.EventPage_Num{
//				Num: 0,
//			},
//			CreatedAt:   &timestamp.Timestamp{},
//			Event:       nil,
//			Synchronous: false,
//		},
//	}
//
//	snapshot := &evented_core.Snapshot{
//		Sequence: 0,
//		State:    nil,
//	}
//
//	book := evented_core.EventBook{
//		Cover:    cover,
//		Pages:    pages,
//		Snapshot: snapshot,
//	}
//
//	ebr := &RepositoryBasic{
//		EventRepo:    erMock.EventRepository{},
//		SnapshotRepo: ssMock.SnapshotRepo{},
//		Domain:       o.domain,
//	}
//
//	err := ebr.Put(context.Background(), &book)
//
//	o.Assert().NoError(err)
//}

func (o *EventBookRepositorySuite) Test_Get() {
	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	ctx := context.Background()
	o.snapshotRepository.On("Get", ctx, o.id).Return(snapshot, nil)
	//ch := make(chan *evented_core.EventPage)
	o.eventRepository.On("GetFrom", ctx, mock.Anything, o.id, uint32(0)).Return(nil)
	book, err := o.eventBookRepository.Get(ctx, o.id)
	o.Assert().NoError(err)
	o.Assert().EqualValues(o.book, book)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(EventBookRepositorySuite))
}
