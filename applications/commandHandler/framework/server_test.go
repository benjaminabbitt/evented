package framework

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/async/mock"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/google/uuid"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type ServerSuite struct {
	suite.Suite
	log     *zap.SugaredLogger
	domainA string
	domainB string
	ctx     context.Context
}

func (s *ServerSuite) SetupTest() {
	s.log = support.Log()
	defer s.log.Sync()
	s.domainA = "testA"
	s.domainB = "testB"
	s.ctx = context.Background()
}

func (s ServerSuite) Test_Handle() {
	eventBookRepo := new(eventBook.MockEventBookRepository)
	holder := new(transport.MockHolder)
	businessClient := new(client.MockClient)
	server := NewServer(eventBookRepo, holder, businessClient, s.log)

	commandBook := s.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)

	eventBookRepo.On("Get", mock2.Anything, id).Return(s.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  s.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessClient.On("Handle", mock2.Anything, contextualCommand).Return(s.produceBusinessResponse(commandBook), nil)
	eventBookRepo.On("Put", mock2.Anything, s.produceBusinessResponse(commandBook)).Return(nil)

	holder.On("GetProjections").Return([]projector.SyncProjectionTransporter{})
	holder.On("GetSaga").Return([]saga.SyncSagaTransporter{})
	holder.On("GetTransports").Return([]async.EventTransporter{})
	server.Handle(context.Background(), commandBook)
	holder.AssertExpectations(s.T())
	businessClient.AssertExpectations(s.T())
	eventBookRepo.AssertExpectations(s.T())
}

func (s ServerSuite) Test_HandleWithTransports() {
	eventBookRepo := new(eventBook.MockEventBookRepository)
	holder := new(transport.MockHolder)
	businessClient := new(client.MockClient)
	server := NewServer(eventBookRepo, holder, businessClient, s.log)

	commandBook := s.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)
	eventBookRepo.On("Get", mock2.Anything, id).Return(s.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  s.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessResponse := s.produceBusinessResponse(commandBook)
	businessClient.On("Handle", mock2.Anything, contextualCommand).Return(businessResponse, nil)
	eventBookRepo.On("Put", mock2.Anything, businessResponse).Return(nil)

	var syncEventPages []*eventedcore.EventPage
	syncEventPages = append(syncEventPages, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	})

	syncEventBook := &eventedcore.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    syncEventPages,
		Snapshot: nil,
	}

	projection := &eventedcore.Projection{
		Cover:      syncEventBook.Cover,
		Projector:  "test",
		Sequence:   0,
		Projection: nil,
	}

	mockProjector := new(projector.MockProjectorClient)
	mockProjector.On("HandleSync", mock2.Anything, syncEventBook).Return(projection, nil)
	holder.On("GetProjections").Return([]projector.SyncProjectionTransporter{mockProjector})

	sagaResult := &eventedcore.EventBook{
		Cover:    nil,
		Pages:    nil,
		Snapshot: nil,
	}

	mockSaga := new(saga.MockSagaClient)
	mockSaga.On("HandleSync", mock2.Anything, syncEventBook).Return(sagaResult, nil)
	holder.On("GetSaga").Return([]saga.SyncSagaTransporter{mockSaga})

	var asyncEventPages []*eventedcore.EventPage
	asyncEventPages = append(asyncEventPages, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	}, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 2},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	asyncEventBook := &eventedcore.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    asyncEventPages,
		Snapshot: nil,
	}

	mockTransport := new(mock.AsyncTransport)
	mockTransport.On("Handle", mock2.Anything, asyncEventBook).Return(nil)
	holder.On("GetTransports").Return([]async.EventTransporter{mockTransport})

	server.Handle(context.Background(), commandBook)
	mockProjector.AssertExpectations(s.T())
	mockSaga.AssertExpectations(s.T())
	mockTransport.AssertExpectations(s.T())
	holder.AssertExpectations(s.T())
	businessClient.AssertExpectations(s.T())
	eventBookRepo.AssertExpectations(s.T())
}

func (s ServerSuite) produceBusinessResponse(commandBook *eventedcore.CommandBook) *eventedcore.EventBook {
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
		Synchronous: true,
	}, &eventedcore.EventPage{
		Sequence:    &eventedcore.EventPage_Num{Num: 2},
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

func (s ServerSuite) produceHistoricalEventBook(commandBook *eventedcore.CommandBook) *eventedcore.EventBook {
	eventPage := NewEmptyEventPage(0, false)
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
