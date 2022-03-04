package eventBook

import (
	"context"
	"github.com/benjaminabbitt/evented/mocks"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

type EventBookRepositorySuite struct {
	suite.Suite
	domain                    string
	id                        uuid.UUID
	eventBookRepository       RepositoryBasic
	snapshotRepository        *mocks.SnapshotStorer
	eventRepository           *mocks.EventStorer
	eventPageRepositoryStream chan *evented.EventPage
}

func (o *EventBookRepositorySuite) SetupTest() {
	id, _ := uuid.NewRandom()
	o.id = id
	o.domain = "test"

	o.eventRepository = &mocks.EventStorer{}
	o.snapshotRepository = &mocks.SnapshotStorer{}
	o.eventPageRepositoryStream = make(chan *evented.EventPage, 10)

	o.eventBookRepository = RepositoryBasic{
		EventRepo:    o.eventRepository,
		SnapshotRepo: o.snapshotRepository,
		Domain:       "test",
	}
}

func (o *EventBookRepositorySuite) Test_Put() {
	id, _ := uuid.NewRandom()
	pid := evented_proto.UUIDToProto(id)
	ctx := context.Background()

	cover := &evented.Cover{
		Domain: "testPut",
		Root:   &pid,
	}

	pages := []*evented.EventPage{
		&evented.EventPage{
			Sequence: &evented.EventPage_Num{
				Num: 0,
			},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
	}

	snapshot := &evented.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	book := evented.EventBook{
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
	snapshot := &evented.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	ctx := context.Background()
	o.snapshotRepository.On("Get", ctx, o.id).Return(snapshot, nil)
	o.eventRepository.On("GetFrom", ctx, mock.Anything, o.id, uint32(0)).Return(nil)
	root := evented_proto.UUIDToProto(o.id)
	expected := evented.EventBook{
		Cover: &evented.Cover{
			Domain: o.domain,
			Root:   &root,
		},
		Pages: []*evented.EventPage{&evented.EventPage{
			Sequence:    &evented.EventPage_Num{Num: 0},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: &evented.Snapshot{},
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
