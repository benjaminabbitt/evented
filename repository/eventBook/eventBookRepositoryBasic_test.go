package eventBook

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	erMock "github.com/benjaminabbitt/evented/repository/events/mock"
	ssMock "github.com/benjaminabbitt/evented/repository/snapshots/mock"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventBookRepositorySuite struct {
	suite.Suite
	domain                    string
	id                        uuid.UUID
	eventBookRepository       RepositoryBasic
	snapshotRepository        *ssMock.SnapshotRepo
	eventRepository           *erMock.EventRepository
	eventPageRepositoryStream chan *evented_core.EventPage
}

func (o *EventBookRepositorySuite) SetupTest() {
	id, _ := uuid.NewRandom()
	o.id = id
	o.domain = "test"

	o.eventRepository = &erMock.EventRepository{}
	o.snapshotRepository = &ssMock.SnapshotRepo{}
	o.eventPageRepositoryStream = make(chan *evented_core.EventPage, 10)

	o.eventBookRepository = RepositoryBasic{
		EventRepo:             o.eventRepository,
		SnapshotRepo:          o.snapshotRepository,
		Domain:                "test",
		EventPageReturnStream: o.eventPageRepositoryStream,
	}
}

func (o *EventBookRepositorySuite) Test_Put() {
	id, _ := uuid.NewRandom()
	pid := evented_proto.UUIDToProto(id)
	ctx := context.Background()

	cover := &evented_core.Cover{
		Domain: "testPut",
		Root:   &pid,
	}

	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence: &evented_core.EventPage_Num{
				Num: 0,
			},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
	}

	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	book := evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}

	ebr := &RepositoryBasic{
		EventRepo:    o.eventRepository,
		SnapshotRepo: o.snapshotRepository,
		Domain:       o.domain,
	}

	o.eventRepository.On("Add", ctx, id, book.Pages).Return(nil)
	o.snapshotRepository.On("Put", ctx, id, book.Snapshot).Return(nil)
	err := ebr.Put(ctx, &book)

	o.eventRepository.AssertExpectations(o.T())
	o.snapshotRepository.AssertExpectations(o.T())
	o.Assert().NoError(err)
}

func (o *EventBookRepositorySuite) Test_Get() {
	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	ctx := context.Background()
	o.snapshotRepository.On("Get", ctx, o.id).Return(snapshot, nil)
	o.eventRepository.On("GetFrom", ctx, o.eventPageRepositoryStream, o.id, uint32(0)).Return(nil)
	root := evented_proto.UUIDToProto(o.id)
	expected := evented_core.EventBook{
		Cover: &evented_core.Cover{
			Domain: o.domain,
			Root:   &root,
		},
		Pages: []*evented_core.EventPage{&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 0},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: &evented_core.Snapshot{},
	}
	o.eventPageRepositoryStream <- expected.Pages[0]
	close(o.eventPageRepositoryStream)
	book, err := o.eventBookRepository.Get(ctx, o.id)
	o.eventRepository.AssertExpectations(o.T())
	o.snapshotRepository.AssertExpectations(o.T())
	o.Assert().NoError(err)
	o.Assert().EqualValues(&expected, book)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(EventBookRepositorySuite))
}
