package framework

import (
	"context"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport"
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type ServerSuite struct {
	suite.Suite
	log     *zap.SugaredLogger
	errh    *evented.ErrLogger
	domainA string
	domainB string
}

func (s *ServerSuite) SetupTest() {
	s.log, s.errh = support.Log()
	defer s.log.Sync()
	s.domainA = "testA"
	s.domainB = "testB"
}

func (s ServerSuite) Test_Handle() {
	eventBookRepo := new(eventBook.MockEventBookRepository)
	holder := new(transport.MockHolder)
	businessClient := new(client.MockClient)
	server := NewServer(eventBookRepo, holder, businessClient, s.log, s.errh)

	commandBook := s.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(*commandBook.Cover.Root)
	eventBookRepo.On("Get", id).Return(*s.produceHistoricalEventBook(*commandBook), nil)

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  s.produceHistoricalEventBook(*commandBook),
		Command: commandBook,
	}

	businessClient.On("Handle", contextualCommand).Return(s.produceBusinessResponse(*commandBook), nil)
	eventBookRepo.On("Put", *s.produceBusinessResponse(*commandBook)).Return(nil)

	holder.On("GetProjections").Return([]projector.SyncProjection{})
	holder.On("GetSaga").Return([]saga.SyncSaga{})
	holder.On("GetTransports").Return([]async.Transport{})
	server.Handle(context.Background(), commandBook)
}

func (s ServerSuite) produceBusinessResponse(commandBook eventedcore.CommandBook) *eventedcore.EventBook {
	var businessReturnEventPages []*eventedcore.EventPage

	businessReturnEventPages = append(businessReturnEventPages, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	businessReturnEventBook := &eventedcore.EventBook{
		Cover:    commandBook.Cover,
		Pages:    businessReturnEventPages,
		Snapshot: nil,
	}

	return businessReturnEventBook
}

func (s ServerSuite) produceHistoricalEventBook(commandBook eventedcore.CommandBook) *eventedcore.EventBook {
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	eventPage := NewEventPage(0, false, *anyEmpty)
	priorStateEventPages := []*eventedcore.EventPage{
		eventPage,
	}

	eb := &eventedcore.EventBook{
		Cover:    commandBook.Cover,
		Pages:    priorStateEventPages,
		Snapshot: nil,
	}
	return eb
}

func (s ServerSuite) produceCommandBook() *eventedcore.CommandBook {
	page := &eventedcore.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     nil,
	}

	randomId, _ := uuid.NewRandom()
	id := evented_proto.UUIDToProto(randomId)

	cover := &eventedcore.Cover{
		Domain: "test",
		Root:   &id,
	}

	commandBook := &eventedcore.CommandBook{
		Cover: cover,
		Pages: []*eventedcore.CommandPage{page},
	}
	return commandBook
}



func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

