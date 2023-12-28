package eventBook

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	mock_events "github.com/benjaminabbitt/evented/repository/events/mocks"
	mock_snapshots "github.com/benjaminabbitt/evented/repository/snapshots/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

type EventBookTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	id     uuid.UUID
	ctx    context.Context
	pid    evented.UUID
	domain string
}

func TestEventBookTestSuite(t *testing.T) {
	suite.Run(t, new(EventBookTestSuite))
}

func (suite *EventBookTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	defer suite.ctrl.Finish()
	suite.id, _ = uuid.NewRandom()
	suite.pid = evented_proto.UUIDToProto(suite.id)
	suite.ctx = context.Background()
	suite.domain = "test"
}

func (suite *EventBookTestSuite) TestPut() {

	eventRepository := mock_events.NewMockEventStorer(suite.ctrl)
	snapshotRepository := mock_snapshots.NewMockSnapshotStorer(suite.ctrl)

	eventBookRepository := RepositoryBasic{
		EventRepo:    eventRepository,
		SnapshotRepo: snapshotRepository,
		Domain:       suite.domain,
	}

	cover := &evented.Cover{
		Domain: "testPut",
		Root:   &suite.pid,
	}

	pages := []*evented.EventPage{
		{
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

	eventRepository.EXPECT().Add(suite.ctx, suite.id, book.Pages).Return(nil)
	snapshotRepository.EXPECT().Put(suite.ctx, suite.id, book.Snapshot).Return(nil)
	err := eventBookRepository.Put(suite.ctx, &book)
	suite.NoError(err)
}

func (suite *EventBookTestSuite) Test_Get() {
	snapshot := &evented.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	root := evented_proto.UUIDToProto(suite.id)
	expected := evented.EventBook{
		Cover: &evented.Cover{
			Domain: suite.domain,
			Root:   &root,
		},
		Pages: []*evented.EventPage{{
			Sequence:    &evented.EventPage_Num{Num: 0},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: &evented.Snapshot{},
	}
	snapshotRepository := mock_snapshots.NewMockSnapshotStorer(suite.ctrl)
	snapshotRepository.EXPECT().Get(suite.ctx, suite.id).Return(snapshot, nil)

	eventRepository := mock_events.NewMockEventStorer(suite.ctrl)
	eventRepository.EXPECT().
		GetFrom(suite.ctx, gomock.Any(), suite.id, uint32(0)).
		Do(func(ctx context.Context, ch chan *evented.EventPage, id uuid.UUID, from uint32) {
			ch <- expected.Pages[0]
			close(ch)
		})

	eventBookRepository := RepositoryBasic{
		EventRepo:    eventRepository,
		SnapshotRepo: snapshotRepository,
		Domain:       suite.domain,
	}

	book, err := eventBookRepository.Get(suite.ctx, suite.id)
	suite.NoError(err)
	suite.EqualValues(&expected, book)
}
