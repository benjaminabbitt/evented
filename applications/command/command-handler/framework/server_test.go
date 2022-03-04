package framework

import (
	"context"
	"errors"
	"github.com/benjaminabbitt/evented/mocks"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	transportMock "github.com/benjaminabbitt/evented/transport/async/mock"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type ServerSuite struct {
	suite.Suite
	log            *zap.SugaredLogger
	domainA        string
	domainB        string
	ctx            context.Context
	eventBookRepo  *mocks.Storer
	holder         *mocks.Holder
	businessClient *mocks.BusinessClient
	server         Server
}

func (o *ServerSuite) SetupTest() {
	o.log = support.Log()
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Error(err)
		}
	}(o.log)
	o.domainA = "testA"
	o.domainB = "testB"
	o.ctx = context.Background()
	o.eventBookRepo = new(mocks.Storer)
	o.holder = new(mocks.Holder)
	o.businessClient = new(mocks.BusinessClient)
	o.server = NewServer(o.eventBookRepo, o.holder, o.businessClient, o.log)
}

func (o ServerSuite) Test_Handle() {
	eventBookRepo := new(mocks.Storer)
	holder := new(mocks.Holder)
	businessClient := new(mocks.BusinessClient)
	server := NewServer(eventBookRepo, holder, businessClient, o.log)

	commandBook := o.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)

	eventBookRepo.On("Get", mock.Anything, id).Return(o.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &evented.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessClient.On("Handle", mock.Anything, contextualCommand).Return(o.produceBusinessResponse(commandBook), nil)
	eventBookRepo.On("Put", mock.Anything, o.produceBusinessResponse(commandBook)).Return(nil)

	holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{})
	holder.On("GetSaga").Return([]saga.SyncSagaTransporter{})
	holder.On("GetTransports").Return([]chan *evented.EventBook{})
	server.Handle(context.Background(), commandBook)
	holder.AssertExpectations(o.T())
	businessClient.AssertExpectations(o.T())
	eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) Test_HandleWithTransports() {

	commandBook := o.produceCommandBook()

	id, _ := evented_proto.ProtoToUUID(commandBook.Cover.Root)
	o.eventBookRepo.On("Get", mock.Anything, id).Return(o.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &evented.ContextualCommand{
		Events:  o.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessResponse := o.produceBusinessResponse(commandBook)
	o.businessClient.On("Handle", mock.Anything, contextualCommand).Return(businessResponse, nil)
	o.eventBookRepo.On("Put", mock.Anything, businessResponse).Return(nil)

	var syncEventPages []*evented.EventPage
	syncEventPages = append(syncEventPages, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	})

	syncEventBook := &evented.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    syncEventPages,
		Snapshot: nil,
	}

	projection := &evented.Projection{
		Cover:      syncEventBook.Cover,
		Projector:  "test",
		Sequence:   0,
		Projection: nil,
	}

	mockProjector := new(mocks.ProjectorClient)
	mockProjector.On("HandleSync", mock.Anything, syncEventBook).Return(projection, nil)
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{mockProjector})

	sagaResult := &evented.SynchronousProcessingResponse{}

	mockSaga := new(mocks.SyncSagaTransporter)
	mockSaga.On("HandleSync", mock.Anything, syncEventBook).Return(sagaResult, nil)
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{mockSaga})

	var asyncEventPages []*evented.EventPage
	asyncEventPages = append(asyncEventPages, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	}, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 2},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	asyncEventBook := &evented.EventBook{
		Cover:    businessResponse.Cover,
		Pages:    asyncEventPages,
		Snapshot: nil,
	}

	mockTransport := new(transportMock.AsyncTransport)
	mockTransport.On("Handle", mock.Anything, asyncEventBook).Return(nil)
	ch := make(chan *evented.EventBook, 10)
	o.holder.On("GetTransports").Return([]chan *evented.EventBook{ch})
	o.server.Handle(context.Background(), commandBook)
	test := <-ch
	o.log.Info(test)
	mockProjector.AssertExpectations(o.T())
	mockSaga.AssertExpectations(o.T())
	o.holder.AssertExpectations(o.T())
	o.businessClient.AssertExpectations(o.T())
	o.eventBookRepo.AssertExpectations(o.T())
}

func (o ServerSuite) produceBusinessResponse(commandBook *evented.CommandBook) *evented.EventBook {
	var businessReturnEventPages []*evented.EventPage

	businessReturnEventPages = append(businessReturnEventPages, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 0},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	}, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 1},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: true,
	}, &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: 2},
		CreatedAt:   nil,
		Event:       nil,
		Synchronous: false,
	})

	businessReturnEventBook := &evented.EventBook{
		Cover:    commandBook.Cover,
		Pages:    businessReturnEventPages,
		Snapshot: nil,
	}

	return businessReturnEventBook
}

func (o ServerSuite) produceHistoricalEventBook(commandBook *evented.CommandBook) *evented.EventBook {
	eventPage := NewEmptyEventPage(0, false)
	priorStateEventPages := []*evented.EventPage{
		eventPage,
	}

	eb := &evented.EventBook{
		Cover:    commandBook.Cover,
		Pages:    priorStateEventPages,
		Snapshot: nil,
	}
	return eb
}

func (o ServerSuite) produceCommandBook() *evented.CommandBook {
	page := &evented.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     nil,
	}

	randomId, _ := uuid.NewRandom()
	id := evented_proto.UUIDToProto(randomId)

	cover := &evented.Cover{
		Domain: "test",
		Root:   &id,
	}

	commandBook := &evented.CommandBook{
		Cover: cover,
		Pages: []*evented.CommandPage{page},
	}
	return commandBook
}

func (o ServerSuite) TestHandleUUIDCorrupt() {
	invalidUUID := &evented.UUID{Value: []byte{}}
	book := o.produceCommandBook()
	book.Cover.Root = invalidUUID
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestEventBookRepositoryError() {
	var typeCheckingBypass *evented.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(typeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestBusinessClientError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}
func (o ServerSuite) TestEventBookRepositoryPutError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(errors.New(""))
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func (o ServerSuite) TestHandleSyncSagaError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*evented.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	sagaTransporter := &mocks.SyncSagaTransporter{}
	sagaTransporter2 := &mocks.SyncSagaTransporter{}
	o.holder.On("GetSaga").Return([]saga.SyncSagaTransporter{sagaTransporter, sagaTransporter2})
	sagaResponse := &evented.SynchronousProcessingResponse{
		Books:       []*evented.EventBook{NewEventBook(id, "", eventPages, nil)},
		Projections: nil,
	}
	sagaTransporter2.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, errors.New(""))
	sagaTransporter.On("HandleSync", mock.Anything, mock.Anything).Return(sagaResponse, nil)
	projectorTransporter := &mocks.ProjectorClient{}
	projectorTransporter2 := &mocks.ProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *evented.Projection = nil
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
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*evented.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", eventPages, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	projectorTransporter := &mocks.ProjectorClient{}
	projectorTransporter2 := &mocks.ProjectorClient{}
	o.holder.On("GetProjectors").Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *evented.Projection = nil
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
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	o.eventBookRepo.On("Get", mock.Anything, mock.Anything).Return(eventBookTypeCheckingBypass, nil)
	id, _ := uuid.NewRandom()
	o.businessClient.On("Handle", mock.Anything, mock.Anything).Return(NewEventBook(id, "", []*evented.EventPage{}, nil), nil)
	o.eventBookRepo.On("Put", mock.Anything, mock.Anything).Return(nil)
	book := o.produceCommandBook()
	_, err := o.server.Handle(o.ctx, book)
	o.Assert().Error(err)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
