package framework

import (
	"context"
	"errors"
	mock_client "github.com/benjaminabbitt/evented/applications/command/command-handler/business/client/mocks"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	mock_transport "github.com/benjaminabbitt/evented/applications/command/command-handler/framework/transport/mocks"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	mock_evented "github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/mocks"
	mock_eventBook "github.com/benjaminabbitt/evented/repository/eventBook/mocks"
	"github.com/benjaminabbitt/evented/support"
	mock_async "github.com/benjaminabbitt/evented/transport/async/mocks"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	mock_projector "github.com/benjaminabbitt/evented/transport/sync/projector/mocks"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	mock_saga "github.com/benjaminabbitt/evented/transport/sync/saga/mocks"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServerSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	actx           *BasicCommandHandlerApplicationContext
	domainA        string
	domainB        string
	ctx            context.Context
	eventBookRepo  *mock_eventBook.MockStorer
	holder         *mock_transport.MockHolder
	businessClient *mock_client.MockBusinessClient
	server         Server
}

func (suite *ServerSuite) SetupTest() {
	log := support.Log()
	suite.ctrl = gomock.NewController(suite.T())
	retryStrategy := &backoff.StopBackOff{}
	config := &configuration.Configuration{}
	suite.actx = &BasicCommandHandlerApplicationContext{
		BasicApplicationContext: support.BasicApplicationContext{
			RetryStrategy: retryStrategy,
			Log:           log,
		},
		Config: config,
	}

	defer func(actx support.ApplicationContext) {
		err := log.Sync()
		if err != nil {
			log.Error(err)
		}
	}(suite.actx)
	suite.domainA = "testA"
	suite.domainB = "testB"
	suite.ctx = context.Background()
	suite.eventBookRepo = mock_eventBook.NewMockStorer(suite.ctrl)
	suite.holder = mock_transport.NewMockHolder(suite.ctrl)
	suite.businessClient = mock_client.NewMockBusinessClient(suite.ctrl)
	suite.server = NewServer(suite.actx, suite.eventBookRepo, suite.holder, suite.businessClient)
}

func (suite ServerSuite) Test_Handle() {
	server := NewServer(suite.actx, suite.eventBookRepo, suite.holder, suite.businessClient)

	commandBook := suite.produceCommandBook()

	id, _ := eventedproto.ProtoToUUID(commandBook.Cover.Root)

	suite.eventBookRepo.EXPECT().
		Get(gomock.Any(), id).
		Return(suite.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &evented.ContextualCommand{
		Events:  suite.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	suite.businessClient.EXPECT().
		Handle(gomock.Any(), contextualCommand).
		Return(suite.produceBusinessResponse(commandBook), nil)

	suite.eventBookRepo.EXPECT().
		Put(gomock.Any(), suite.produceBusinessResponse(commandBook)).
		Return(nil)

	suite.holder.EXPECT().
		GetProjectors().
		Return([]projector.SyncProjectorTransporter{})

	suite.holder.EXPECT().
		GetSaga().
		Return([]saga.SyncSagaTransporter{})

	suite.holder.EXPECT().
		GetTransports().
		Return([]chan *evented.EventBook{})

	_, err := server.Handle(context.Background(), commandBook)
	assert.NoError(suite.T(), err)
}

func (suite ServerSuite) Test_HandleWithTransports() {

	commandBook := suite.produceCommandBook()

	id, _ := eventedproto.ProtoToUUID(commandBook.Cover.Root)

	suite.eventBookRepo.EXPECT().Get(gomock.Any(), id).
		Return(suite.produceHistoricalEventBook(commandBook), nil)

	contextualCommand := &evented.ContextualCommand{
		Events:  suite.produceHistoricalEventBook(commandBook),
		Command: commandBook,
	}

	businessResponse := suite.produceBusinessResponse(commandBook)
	suite.businessClient.EXPECT().
		Handle(gomock.Any(), contextualCommand).
		Return(businessResponse, nil)

	suite.eventBookRepo.EXPECT().
		Put(gomock.Any(), businessResponse).
		Return(nil)

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

	mockProjector := mock_evented.NewMockProjectorClient(suite.ctrl)
	mockProjector.EXPECT().
		HandleSync(gomock.Any(), syncEventBook).
		Return(projection, nil)

	suite.holder.EXPECT().
		GetProjectors().
		Return([]projector.SyncProjectorTransporter{mockProjector})

	sagaResult := &evented.SynchronousProcessingResponse{}

	mockSaga := mock_evented.NewMockSagaClient(suite.ctrl)
	mockSaga.EXPECT().HandleSync(gomock.Any(), syncEventBook).Return(sagaResult, nil)
	suite.holder.EXPECT().GetSaga().Return([]saga.SyncSagaTransporter{mockSaga})

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

	transporter := mock_async.NewMockEventTransporter(suite.ctrl)
	transporter.EXPECT().Handle(gomock.Any(), asyncEventBook).Return(nil)
	ch := make(chan *evented.EventBook, 10)
	suite.holder.EXPECT().GetTransports().Return([]chan *evented.EventBook{ch})
	_, err := suite.server.Handle(context.Background(), commandBook)
	if err != nil {
		suite.Error(err)
	}
}

func (suite ServerSuite) produceBusinessResponse(commandBook *evented.CommandBook) *evented.EventBook {
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

func (suite ServerSuite) produceHistoricalEventBook(commandBook *evented.CommandBook) *evented.EventBook {
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

func (suite ServerSuite) produceCommandBook() *evented.CommandBook {
	page := &evented.CommandPage{
		Sequence:    0,
		Synchronous: false,
		Command:     nil,
	}

	randomId, _ := uuid.NewRandom()
	id := eventedproto.UUIDToProto(randomId)

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

func (suite ServerSuite) TestHandleUUIDCorrupt() {
	invalidUUID := &evented.UUID{Value: []byte{}}
	book := suite.produceCommandBook()
	book.Cover.Root = invalidUUID
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestEventBookRepositoryError() {
	var typeCheckingBypass *evented.EventBook = nil
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(typeCheckingBypass, errors.New(""))
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestBusinessClientError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	suite.businessClient.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, errors.New(""))
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestEventBookRepositoryPutError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	suite.businessClient.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	suite.eventBookRepo.EXPECT().Put(gomock.Any(), gomock.Any()).Return(errors.New(""))
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestHandleSyncSagaError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*evented.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	suite.businessClient.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(NewEventBook(id, "", eventPages, nil), nil)
	suite.eventBookRepo.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
	sagaTransporter := mock_saga.NewMockSyncSagaTransporter(suite.ctrl)
	sagaTransporter2 := mock_saga.NewMockSyncSagaTransporter(suite.ctrl)
	suite.holder.EXPECT().GetSaga().Return([]saga.SyncSagaTransporter{sagaTransporter, sagaTransporter2})
	sagaResponse := &evented.SynchronousProcessingResponse{
		Books:       []*evented.EventBook{NewEventBook(id, "", eventPages, nil)},
		Projections: nil,
	}
	sagaTransporter.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(sagaResponse, nil)
	sagaTransporter2.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(sagaResponse, errors.New(""))
	projectorTransporter := mock_projector.NewMockSyncProjectorTransporter(suite.ctrl)
	projectorTransporter2 := mock_projector.NewMockSyncProjectorTransporter(suite.ctrl)
	suite.holder.EXPECT().GetProjectors().Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *evented.Projection = nil
	projectorTransporter.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(projectionTypeCheckingBypass, nil)
	projectorTransporter2.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(projectionTypeCheckingBypass, nil)
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestHandleSyncProjectionError() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	id, _ := uuid.NewRandom()
	eventPages := []*evented.EventPage{NewEmptyEventPage(0, false), NewEmptyEventPage(1, false)}
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	suite.businessClient.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(NewEventBook(id, "", eventPages, nil), nil)
	suite.eventBookRepo.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
	projectorTransporter := mock_evented.NewMockProjectorClient(suite.ctrl)
	projectorTransporter2 := mock_evented.NewMockProjectorClient(suite.ctrl)
	suite.holder.EXPECT().GetProjectors().Return([]projector.SyncProjectorTransporter{projectorTransporter, projectorTransporter2})
	var projectionTypeCheckingBypass *evented.Projection = nil
	projectorTransporter.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(projectionTypeCheckingBypass, nil)
	projectorTransporter2.EXPECT().HandleSync(gomock.Any(), gomock.Any()).Return(projectionTypeCheckingBypass, errors.New(""))
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func (suite ServerSuite) TestExtractSynchronousEmptyEventBook() {
	var eventBookTypeCheckingBypass *evented.EventBook = nil
	suite.eventBookRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(eventBookTypeCheckingBypass, nil)
	id, _ := uuid.NewRandom()
	suite.businessClient.EXPECT().Handle(gomock.Any(), gomock.Any()).Return(NewEventBook(id, "", []*evented.EventPage{}, nil), nil)
	suite.eventBookRepo.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
	book := suite.produceCommandBook()
	_, err := suite.server.Handle(suite.ctx, book)
	suite.Assert().Error(err)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
