package eventQueryServer

import (
	"context"
	"crypto/rand"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/query"
	mer "github.com/benjaminabbitt/evented/repository/events/mock"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type QueryHandlerSuite struct {
	suite.Suite
	log   *zap.SugaredLogger
	ctx   context.Context
	repos *mer.EventRepository
	sut   *DefaultEventQueryServer
}

func (o *QueryHandlerSuite) SetupTest() {
	o.log = support.Log()
	defer o.log.Sync()
	o.ctx = context.Background()
	o.repos = &mer.EventRepository{}
	o.sut = &DefaultEventQueryServer{
		eventRepos: o.repos,
		log:        o.log,
	}
}

func sumPages(pages []*eventedcore.EventPage) uint32 {
	size := uint32(0)
	for _, page := range pages {
		size += uint32(proto.Size(page))
	}
	return size
}

/// Validates the memory approximation technique we use when batching events into event books
func (o *QueryHandlerSuite) TestMemoryApproximation() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	cover := &eventedcore.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}

	book := &eventedcore.EventBook{
		Cover:    cover,
		Pages:    nil,
		Snapshot: nil,
	}

	e1 := make([]byte, 16)
	rand.Read(e1)
	pba := &eventedquery.TestByteArray{
		Bytes: e1,
	}
	a1, _ := ptypes.MarshalAny(pba)

	e2 := make([]byte, 32)
	rand.Read(e2)
	pba2 := &eventedquery.TestByteArray{
		Bytes: e2,
	}
	a2, _ := ptypes.MarshalAny(pba2)

	pages := []*eventedcore.EventPage{
		&eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Num{Num: 1},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       a1,
			Synchronous: false,
		},
		&eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Num{Num: 2},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       a2,
			Synchronous: false,
		},
	}

	book2 := &eventedcore.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	addition := uint32(proto.Size(book)) + sumPages(pages)
	composed := proto.Size(book2)
	o.log.Info("addition (bytes): ", addition)
	o.log.Info("composite (bytes): ", proto.Size(book2))
	//Allow up to 1k variance
	o.Assert().InDelta(composed, addition, 1000)
}

func (o *QueryHandlerSuite) Test_Low_High() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	domain := "test"

	//evtChan := make(chan *evented_core.EventPage)
	query := &eventedquery.Query{
		Domain:     domain,
		Root:       &protoUUID,
		LowerBound: 1,
		UpperBound: 2,
	}
	ctx := context.Background()
	o.repos.On("GetFromTo", mock.Anything, mock.Anything, uuid, uint32(1), uint32(2)).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *eventedcore.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*eventedcore.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*eventedcore.EventPage_Num).Num)
		o.Assert().Equal(domain, book.Cover.Domain)
		o.Assert().Equal(&protoUUID, book.Cover.Root)
	})

	_ = o.sut.GetEvents(query, queryResponse)

	o.repos.AssertExpectations(o.T())
	queryResponse.AssertExpectations(o.T())
}

func (o *QueryHandlerSuite) Test_Low() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	domain := "test"

	query := &eventedquery.Query{
		Domain:     domain,
		Root:       &protoUUID,
		LowerBound: 1,
	}
	ctx := context.Background()
	o.repos.On("GetFrom", mock.Anything, mock.Anything, uuid, uint32(1)).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *eventedcore.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*eventedcore.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*eventedcore.EventPage_Num).Num)
		o.Assert().Equal(domain, book.Cover.Domain)
		o.Assert().Equal(&protoUUID, book.Cover.Root)
	})

	_ = o.sut.GetEvents(query, queryResponse)

	o.repos.AssertExpectations(o.T())
	queryResponse.AssertExpectations(o.T())
}

func (o *QueryHandlerSuite) Test_NoLimits() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := eventedproto.UUIDToProto(uuid)
	domain := "test"

	query := &eventedquery.Query{
		Domain: domain,
		Root:   &protoUUID,
	}
	ctx := context.Background()
	o.repos.On("Get", mock.Anything, mock.Anything, uuid).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *eventedcore.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*eventedcore.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*eventedcore.EventPage_Num).Num)
		o.Assert().Equal(domain, book.Cover.Domain)
		o.Assert().Equal(&protoUUID, book.Cover.Root)
	})

	_ = o.sut.GetEvents(query, queryResponse)

	o.repos.AssertExpectations(o.T())
	queryResponse.AssertExpectations(o.T())

}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(QueryHandlerSuite))
}
