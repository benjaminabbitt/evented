package eventBook

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	mock_events "github.com/benjaminabbitt/evented/repository/events/mocks"
	mock_snapshots "github.com/benjaminabbitt/evented/repository/snapshots/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id, _ := uuid.NewRandom()
	pid := evented_proto.UUIDToProto(id)
	ctx := context.Background()
	domain := "test"

	eventRepository := mock_events.NewMockEventStorer(ctrl)
	snapshotRepository := mock_snapshots.NewMockSnapshotStorer(ctrl)

	eventBookRepository := RepositoryBasic{
		EventRepo:    eventRepository,
		SnapshotRepo: snapshotRepository,
		Domain:       domain,
	}

	cover := &evented.Cover{
		Domain: "testPut",
		Root:   &pid,
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

	eventRepository.EXPECT().Add(ctx, id, book.Pages).Return(nil)
	snapshotRepository.EXPECT().Put(ctx, id, book.Snapshot).Return(nil)
	err := eventBookRepository.Put(ctx, &book)
	assert.NoError(t, err)
}

func Test_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id, _ := uuid.NewRandom()
	ctx := context.Background()
	domain := "test"

	snapshot := &evented.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	root := evented_proto.UUIDToProto(id)
	expected := evented.EventBook{
		Cover: &evented.Cover{
			Domain: domain,
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
	snapshotRepository := mock_snapshots.NewMockSnapshotStorer(ctrl)
	snapshotRepository.EXPECT().Get(ctx, id).Return(snapshot, nil)

	eventRepository := mock_events.NewMockEventStorer(ctrl)
	eventRepository.EXPECT().
		GetFrom(ctx, gomock.Any(), id, uint32(0)).
		Do(func(ctx context.Context, ch chan *evented.EventPage, id uuid.UUID, from uint32) {
			ch <- expected.Pages[0]
			close(ch)
		})

	eventBookRepository := RepositoryBasic{
		EventRepo:    eventRepository,
		SnapshotRepo: snapshotRepository,
		Domain:       domain,
	}

	book, err := eventBookRepository.Get(ctx, id)
	assert.NoError(t, err)
	assert.EqualValues(t, &expected, book)
}
