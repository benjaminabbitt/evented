package framework

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	event_memory "github.com/benjaminabbitt/evented/repository/events/event-memory"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/benjaminabbitt/evented/transport"
	projectormock "github.com/benjaminabbitt/evented/transport/sync/projector/mock"
	sagamock "github.com/benjaminabbitt/evented/transport/sync/saga/mock"
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
	eventBookRepo := eventBook.Repository{
		EventRepo:    event_memory.NewMemoryRepository(),
		SnapshotRepo: snapshot_memory.NewSSMemoryRepository(),
	}

	syncSagas := []transport.SyncSaga{sagamock.NewSagaClient()}

	syncProjections := []transport.SyncProjection{projectormock.NewProjectorClient()}

	sagas := []transport.Saga{sagamock.NewSagaClient()}

	projections := []transport.Projection{projectormock.NewProjectorClient()}

	business := &businessLogic.MockBusinessLogicClient{}

	domain := "test"

	server := NewServer(eventBookRepo, syncSagas, syncProjections, sagas, projections, business)

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
		in.Command.Cover.Root,
		in.Command.Cover.Domain,
		[]*evented_core.EventPage{
			NewEventPage(0, false, *anyEmpty),
			NewEventPage(1, true, *anyEmpty),
			NewEventPage(2, false, *anyEmpty),
		},
		nil,
		), nil
}
