package eventQueryServer

import (
	"context"
	"crypto/rand"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	mock_evented "github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/mocks"
	mock_events "github.com/benjaminabbitt/evented/repository/events/mocks"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"

	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/proto"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

type QueryHandlerSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	log           *zap.SugaredLogger
	ctx           context.Context
	repos         *mock_events.MockEventStorer
	sut           *DefaultEventQueryServer
	eventPageChan chan *evented.EventPage
}

func (suite *QueryHandlerSuite) SetupTest() {
	suite.log = support.Log()
	suite.ctrl = gomock.NewController(suite.T())
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(suite.log)
	suite.ctx = context.Background()
	suite.repos = mock_events.NewMockEventStorer(suite.ctrl)
	suite.sut = &DefaultEventQueryServer{
		eventRepos: suite.repos,
		log:        suite.log,
	}
	suite.eventPageChan = make(chan *evented.EventPage)
}

func sumPages(pages []*evented.EventPage) uint32 {
	size := uint32(0)
	for _, page := range pages {
		size += uint32(proto.Size(page))
	}
	return size
}

/// Validates the memory approximation technique we use when batching events into event books
func (suite *QueryHandlerSuite) TestMemoryApproximation() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	cover := &evented.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}

	book := &evented.EventBook{
		Cover:    cover,
		Pages:    nil,
		Snapshot: nil,
	}

	e1 := make([]byte, 16)
	rand.Read(e1)
	pba := &evented.TestByteArray{
		Bytes: e1,
	}
	a1, _ := anypb.New(pba)

	e2 := make([]byte, 32)
	rand.Read(e2)
	pba2 := &evented.TestByteArray{
		Bytes: e2,
	}
	a2, _ := anypb.New(pba2)

	pages := []*evented.EventPage{
		{
			Sequence:    &evented.EventPage_Num{Num: 1},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       a1,
			Synchronous: false,
		},
		{
			Sequence:    &evented.EventPage_Num{Num: 2},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       a2,
			Synchronous: false,
		},
	}

	book2 := &evented.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	addition := uint32(proto.Size(book)) + sumPages(pages)
	composed := proto.Size(book2)
	suite.log.Info("addition (bytes): ", addition)
	suite.log.Info("composite (bytes): ", proto.Size(book2))
	//Allow up to 1k variance
	suite.Assert().InDelta(composed, addition, 1000)
}

func (suite *QueryHandlerSuite) Test_Low_High() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	domain := "test"

	//evtChan := make(chan *evented.EventPage)
	query := &evented.Query{
		Domain:     domain,
		Root:       &protoUUID,
		LowerBound: 1,
		UpperBound: 2,
	}
	ctx := context.Background()
	page := framework.NewEmptyEventPage(1, false)
	suite.repos.EXPECT().
		GetFromTo(gomock.Any(), gomock.Any(), uuid, uint32(1), uint32(2)).
		Do(func() {
			suite.eventPageChan <- page
			close(suite.eventPageChan)
		}).
		Return(nil)

	queryResponse := mock_evented.NewMockEventQuery_GetEventsServer(suite.ctrl)
	queryResponse.EXPECT().Context().Return(ctx)
	queryResponse.EXPECT().
		Send(gomock.Any()).
		Return(nil).
		Do(func(payload interface{}) {
			book := payload.(*evented.EventBook)
			suite.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
			suite.Assert().Equal(domain, book.Cover.Domain)
			suite.Assert().Equal(&protoUUID, book.Cover.Root)
		})

	_ = suite.sut.GetEvents(query, queryResponse)
}

func (suite *QueryHandlerSuite) Test_Low() {
	id, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(id)
	domain := "test"

	query := &evented.Query{
		Domain:     domain,
		Root:       &protoUUID,
		LowerBound: 1,
	}
	ctx := context.Background()
	suite.repos.EXPECT().GetFrom(gomock.Any(), gomock.Any(), id, uint32(1)).
		Do(func(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32) {
			page := framework.NewEmptyEventPage(1, false)
			go func() {
				evtChan <- page
				close(evtChan)
			}()
		}).
		Return(nil)

	queryResponse := mock_evented.NewMockEventQuery_GetEventsServer(suite.ctrl)
	queryResponse.EXPECT().Context().Return(ctx)
	queryResponse.EXPECT().
		Send(gomock.Any()).
		Return(nil).
		Do(func(payload interface{}) {
			book := payload.(*evented.EventBook)
			suite.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
			suite.Assert().Equal(domain, book.Cover.Domain)
			suite.Assert().Equal(&protoUUID, book.Cover.Root)
		})

	_ = suite.sut.GetEvents(query, queryResponse)
}

func (suite *QueryHandlerSuite) Test_NoLimits() {
	id, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(id)
	domain := "test"

	query := &evented.Query{
		Domain: domain,
		Root:   &protoUUID,
	}
	ctx := context.Background()
	suite.repos.EXPECT().
		Get(gomock.Any(), gomock.Any(), id).
		Do(func(ctx context.Context, evtChan chan *evented.EventPage, id uuid.UUID, from uint32) {
			page := framework.NewEmptyEventPage(1, false)
			evtChan <- page
			close(evtChan)
		}).
		Return(nil)

	queryResponse := mock_evented.NewMockEventQuery_GetEventsServer(suite.ctrl)
	queryResponse.EXPECT().Context().Return(ctx)
	queryResponse.EXPECT().
		Send(gomock.Any()).
		Return(nil).
		Do(func(payload interface{}) {
			book := payload.(*evented.EventBook)
			suite.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
			suite.Assert().Equal(domain, book.Cover.Domain)
			suite.Assert().Equal(&protoUUID, book.Cover.Root)
		})

	_ = suite.sut.GetEvents(query, queryResponse)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(QueryHandlerSuite))
}
