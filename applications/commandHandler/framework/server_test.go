package framework

import (
	"context"
	"errors"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
	transportMock "github.com/benjaminabbitt/evented/transport/async/mock"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type ServerSuite struct {
	suite.Suite
	log            *zap.SugaredLogger
	domainA        string
	domainB        string
	ctx            context.Context
	eventBookRepo  *eventBook.MockEventBookRepository
	holder         *transport.MockHolder
	businessClient *client.MockClient
	server         Server
}

func (o *ServerSuite) SetupTest() {
	o.log = support.Log()
	defer o.log.Sync()
	o.domainA = "testA"
	o.domainB = "testB"
	o.ctx = context.Background()
	o.eventBookRepo = new(eventBook.MockEventBookRepository)
	o.holder = new(transport.MockHolder)
	o.businessClient = new(client.MockClient)
	o.server = NewServer(o.eventBookRepo, o.holder, o.businessClient, o.log)
}

func (o ServerSuite) Test_Handle() {
	eventBookRepo := new(eventBook.MockEventBookRepository)
	holder := new(transport.MockHolder)
	businessClient := new(client.MockClient)
	server := NewServer(eventBookRepo, holder, businessClient, o.log)

	commandBook := o.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)

	eventBookRepo.On("Get", mock.Anything, id).Return(o.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessClient.On("Handle", mock.Anything, contextualCommand).Return(o.produceBusinessResponse(commandBook), nil)
	eventBookRepo.On("Put", mock.Anything, o.produceBusinessResponse(commandBook)).Return(nil)

	holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{})
	holder.On("GetSaga").Return([]saga.SyncSagaTransporter{})
	holder.On("GetTransports").Return([]chan *eventedcore.EventBook{})
	server.Handle(context.Background(), commandBook)
	holder.AssertExpectations(o.T())
	businessClient.AssertExpectations(o.T())
	eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) Test_HandleWithTransports() {

	commandBook := o.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)
	o.eventBookRepo.On("Get", mock.Anything, id).Return(o.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessResponse := o.produceBusinessResponse(commandBook)
	o.businessClient.On("Handle", mock.Anything, contextualCommand).Return(businessResponse, nil)
	o.eventBookRepo.On("Put", mock.Anything, businessResponse).Return(nil)

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
	mockProjector.On("HandleSync", mock.Anything, syncEventBook).Return(projection, nil)
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{mockProjector})

	sagaResult := &eventedcore.SynchronousProcessingResponse{}

	mockSaga := new(saga.MockSagaClient)
	mockSaga.On("HandleSync", mock.Anything, syncEventBook).Return(sagaResult, nil)
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{mockSaga})

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

	mockTransport := new(transportMock.AsyncTransport)
	mockTransport.On("Handle", mock.Anything, asyncEventBook).Return(nil)
	ch := make(chan *eventedcore.EventBook, 10)
	o.holder.On("GetTransports").Return([]chan *eventedcore.EventBook{ch})
	o.server.Handle(context.Background(), commandBook)
	test := <-ch
	o.log.Info(test)
	mockProjector.AssertExpectations(o.T())
	mockSaga.AssertExpectations(o.T())
	o.holder.AssertExpectations(o.T())
	o.businessClient.AssertExpectations(o.T())
	o.eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) produceBusinessResponse(commandBook *eventedcore.CommandBook) *eventedcore.EventBook {
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

func (o ServerSuite) produceHistoricalEventBook(commandBook *eventedcore.CommandBook) *eventedcore.EventBook {
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

func (o ServerSuite) produceCommandBook() *eventedcore.CommandBook {
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

func (o ServerSuite) TestHandleUUIDCorrupt() {
	invalidUUID := &eventedcore.UUID{Value: []byte{}}
	book := o.produceCommandBook()
	book.Cover.Root = invalidUUID
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestEventBookRepositoryError() {
	var typeCheckingBypass *eventedcore.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(typeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestBusinessClientError() {
	var eventBookTypeCheckingBypass *eventedcore.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}
func (o ServerSuite) TestEventBookRepositoryPutError() {
	var eventBookTypeCheckingBypass *eventedcore.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestHandleSyncSagaError() {
	var eventBookTypeCheckingBypass *eventedcore.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*eventedcore.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	sagaTransporter := &saga.MockSagaClient{}
	sagaTransporter2 := &saga.MockSagaClient{}
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{sagaTransporter, sagaTransporter2})
	sagaResponse := &eventedcore.SynchronousProcessingResponse{
		Books:       []*eventedcore.EventBook{NewEventBook(id, "", eventPages, nil)},
		Projections: nil,
	}
	sagaTransporter2.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, errors.New(""))
	sagaTransporter.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, nil)
	projectorTransporter := &projector.MockProjectorClient{}
	projectorTransporter2 := &projector.MockProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *eventedcore.Projection = nil
	projectorTransporter.On("HandleSync", mock.Anything, mock.Anything).Return(projectionTypeCheckingBypass, nil)
	projectorTransporter2.On("HandleSync", mock.Anything, mock.Anything).Return(projectionTypeCheckingBypass, nil)
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
	sagaTransporter.AssertExpectations(o.T())
	sagaTransporter2.AssertExpectations(o.T())
	projectorTransporter.AssertExpectations(o.T())
	projectorTransporter2.AssertExpectations(o.T())
	o.holder.AssertExpectations(o.T())
	o.businessClient.AssertExpectations(o.T())
	o.eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) TestHandleSyncProjectionError() {
	var eventBookTypeCheckingBypass *eventedcore.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*eventedcore.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	projectorTransporter := &projector.MockProjectorClient{}
	projectorTransporter2 := &projector.MockProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *eventedcore.Projection = nil
	projectorTransporter.On("HandleSync", mock.Anything, mock.Anything).Return(projectionTypeCheckingBypass, nil)
	projectorTransporter2.On("HandleSync", mock.Anything, mock.Anything).Return(projectionTypeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
	projectorTransporter.AssertExpectations(o.T())
	projectorTransporter2.AssertExpectations(o.T())
	o.holder.AssertExpectations(o.T())
	o.businessClient.AssertExpectations(o.T())
	o.eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) TestExtractSynchronousEmptyEventBook() {
	var eventBookTypeCheckingBypass *eventedcore.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	id, _ := uuid.NewRandom()
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", []*eventedcore.EventPage{}, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestListenForNoErrors() {
	var err error
	defer o.server.Earmuffs()
	go func() {
		err = o.server.Listen(1000)
	}()
	time.Sleep(1 * time.Second)
	o.Assert().NoError(err)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
