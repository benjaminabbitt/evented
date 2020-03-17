package framework

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	event_memory "github.com/benjaminabbitt/evented/repository/events/event-memory"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	async2 "github.com/benjaminabbitt/evented/transport/async"
	async "github.com/benjaminabbitt/evented/transport/projector/console"
	sync "github.com/benjaminabbitt/evented/transport/sync/projector/console"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"testing"
)

import "github.com/google/uuid"

type ServerSuite struct {
	suite.Suite
}

func (s *ServerSuite) Test_Handle(){
	ctx := context.Background()
	eventRepo := event_memory.NewMemoryRepository()
	ssRepo := snapshot_memory.NewSSMemoryRepository()
	asyncSenders := []async2.EventSender{&async.Sender{}}
	syncSender := []saga.SyncSaga{&sync.Sender{}}

	domain := "test"

	server := NewServer(eventRepo, ssRepo, asyncSenders, syncSender, &BusinessLogicMock{})
	id, _ := uuid.NewRandom()
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	page := &evented_core.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     anyEmpty,
	}
	commandBook := &evented_core.CommandBook{
		Cover: &evented_core.Cover{
			Domain: domain,
			Root:     id.String(),
		},
		Pages: []*evented_core.CommandPage{page},
	}
	commandResponse, _ := server.Handle(ctx, commandBook)

	s.Assert().EqualValues(0, commandResponse.Books[0].Pages[0].Sequence)
	s.Assert().Equal(false, commandResponse.Books[0].Pages[0].Synchronous)
	s.Assert().EqualValues(1, commandResponse.Books[0].Pages[1].Sequence)
	s.Assert().Equal(true, commandResponse.Books[0].Pages[1].Synchronous)
	s.Assert().Equal(id.String(), commandResponse.Books[0].Cover.Root)
	s.Assert().Equal(domain, commandResponse.Books[0].Cover.Domain)
}


func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

type BusinessLogicMock struct{}

func (bl *BusinessLogicMock) Handle(ctx context.Context, in *evented_core.ContextualCommand, opts ...grpc.CallOption) (*evented_core.EventBook, error){
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	return NewEventBook(
		in.Command.Cover.Id,
		in.Command.Cover.Domain,
		[]*evented_core.EventPage{
			NewEventPage(0, false, *anyEmpty),
			NewEventPage(1, true, *anyEmpty),
			NewEventPage(2, false, *anyEmpty),
		},
		nil,
		), nil
}
