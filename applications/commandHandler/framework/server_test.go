package framework

import (
	"context"
	"errors"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
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

	contextualCommand := &core.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessClient.On("Handle", mock.Anything, contextualCommand).Return(o.produceBusinessResponse(commandBook), nil)
	eventBookRepo.On("Put", mock.Anything, o.produceBusinessResponse(commandBook)).Return(nil)

	holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{})
	holder.On("GetSaga").Return([]saga.SyncSagaTransporter{})
	holder.On("GetTransports").Return([]chan *core.EventBook{})
	server.Handle(context.Background(), commandBook)
	holder.AssertExpectations(o.T())
	businessClient.AssertExpectations(o.T())
	eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) Test_HandleWithTransports() {

	commandBook := o.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)
	o.eventBookRepo.On("Get", mock.Anything, id).Return(o.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &core.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessResponse := o.produceBusinessResponse(commandBook)
	o.businessClient.On("Handle", mock.Anything, contextualCommand).Return(businessResponse, nil)
	o.eventBookRepo.On("Put", mock.Anything, businessResponse).Return(nil)

	var syncEventPages []*core.EventPage
	syncEventPages = append(syncEventPages, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	})

	syncEventBook := &core.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    syncEventPages,
		Snapshot: nil,
	}

	projection := &core.Projection{
		Cover:      syncEventBook.Cover,
		Projector:  "test",
		Sequence:   0,
		Projection: nil,
	}

	mockProjector := new(projector.MockProjectorClient)
	mockProjector.On("HandleSync", mock.Anything, syncEventBook).Return(projection, nil)
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{mockProjector})

	sagaResult := &core.SynchronousProcessingResponse{}

	mockSaga := new(saga.MockSagaClient)
	mockSaga.On("HandleSync", mock.Anything, syncEventBook).Return(sagaResult, nil)
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{mockSaga})

	var asyncEventPages []*core.EventPage
	asyncEventPages = append(asyncEventPages, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	}, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 2},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	asyncEventBook := &core.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    asyncEventPages,
		Snapshot: nil,
	}

	mockTransport := new(transportMock.AsyncTransport)
	mockTransport.On("Handle", mock.Anything, asyncEventBook).Return(nil)
	ch := make(chan *core.EventBook, 10)
	o.holder.On("GetTransports").Return([]chan *core.EventBook{ch})
	o.server.Handle(context.Background(), commandBook)
	test := <-ch
	o.log.Info(test)
	mockProjector.AssertExpectations(o.T())
	mockSaga.AssertExpectations(o.T())
	o.holder.AssertExpectations(o.T())
	o.businessClient.AssertExpectations(o.T())
	o.eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) produceBusinessResponse(commandBook *core.CommandBook) *core.EventBook {
	var businessReturnEventPages []*core.EventPage

	businessReturnEventPages = append(businessReturnEventPages, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	}, &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 2},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	businessReturnEventBook := &core.EventBook{
		Cover:    commandBook.Cover,
		Pages:    businessReturnEventPages,
		Snapshot: nil,
	}

	return businessReturnEventBook
}

func (o ServerSuite) produceHistoricalEventBook(commandBook *core.CommandBook) *core.EventBook {
	eventPage := NewEmptyEventPage(0, false)
	priorStateEventPages := []*core.EventPage{
		eventPage,
	}

	eb := &core.EventBook{
		Cover:    commandBook.Cover,
		Pages:    priorStateEventPages,
		Snapshot: nil,
	}
	return eb
}

func (o ServerSuite) produceCommandBook() *core.CommandBook {
	page := &core.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     nil,
	}

	randomId, _ := uuid.NewRandom()
	id := evented_proto.UUIDToProto(randomId)

	cover := &core.Cover{
		Domain: "test",
		Root:   &id,
	}

	commandBook := &core.CommandBook{
		Cover: cover,
		Pages: []*core.CommandPage{page},
	}
	return commandBook
}

func (o ServerSuite) TestHandleUUIDCorrupt() {
	invalidUUID := &core.UUID{Value: []byte{}}
	book := o.produceCommandBook()
	book.Cover.Root = invalidUUID
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestEventBookRepositoryError() {
	var typeCheckingBypass *core.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(typeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestBusinessClientError() {
	var eventBookTypeCheckingBypass *core.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}
func (o ServerSuite) TestEventBookRepositoryPutError() {
	var eventBookTypeCheckingBypass *core.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestHandleSyncSagaError() {
	var eventBookTypeCheckingBypass *core.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*core.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	sagaTransporter := &saga.MockSagaClient{}
	sagaTransporter2 := &saga.MockSagaClient{}
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{sagaTransporter, sagaTransporter2})
	sagaResponse := &core.SynchronousProcessingResponse{
		Books:       []*core.EventBook{NewEventBook(id, "", eventPages, nil)},
		Projections: nil,
	}
	sagaTransporter2.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, errors.New(""))
	sagaTransporter.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, nil)
	projectorTransporter := &projector.MockProjectorClient{}
	projectorTransporter2 := &projector.MockProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *core.Projection = nil
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
	var eventBookTypeCheckingBypass *core.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*core.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	projectorTransporter := &projector.MockProjectorClient{}
	projectorTransporter2 := &projector.MockProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *core.Projection = nil
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
	var eventBookTypeCheckingBypass *core.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	id, _ := uuid.NewRandom()
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", []*core.EventPage{}, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestListenForNoErrors() {
	var err error
	defer o.server.Shutdown()
	go func() {
		err = o.server.Listen(1000)
	}()
	time.Sleep(1 * time.Second)
	o.Assert().NoError(err)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
