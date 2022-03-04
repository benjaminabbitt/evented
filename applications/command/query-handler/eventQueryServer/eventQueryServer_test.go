package eventQueryServer

import (
	"context"
	"crypto/rand"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	mer "github.com/benjaminabbitt/evented/repository/events/mock"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/proto"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(o.log)
	o.ctx = context.Background()
	o.repos = &mer.EventRepository{}
	o.sut = &DefaultEventQueryServer{
		eventRepos: o.repos,
		log:        o.log,
	}
}

func sumPages(pages []*evented.EventPage) uint32 {
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
	o.log.Info("addition (bytes): ", addition)
	o.log.Info("composite (bytes): ", proto.Size(book2))
	//Allow up to 1k variance
	o.Assert().InDelta(composed, addition, 1000)
}

func (o *QueryHandlerSuite) Test_Low_High() {
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
	o.repos.On("GetFromTo", mock.Anything, mock.Anything, uuid, uint32(1), uint32(2)).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *evented.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := &MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*evented.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
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

	query := &evented.Query{
		Domain:     domain,
		Root:       &protoUUID,
		LowerBound: 1,
	}
	ctx := context.Background()
	o.repos.On("GetFrom", mock.Anything, mock.Anything, uuid, uint32(1)).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *evented.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := &MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*evented.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
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

	query := &evented.Query{
		Domain: domain,
		Root:   &protoUUID,
	}
	ctx := context.Background()
	o.repos.On("Get", mock.Anything, mock.Anything, uuid).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *evented.EventPage)
		page := framework.NewEmptyEventPage(1, false)
		go func() {
			ch <- page
			close(ch)
		}()
	}).Once()

	queryResponse := &MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		book := args.Get(0).(*evented.EventBook)
		o.Assert().Equal(uint32(1), book.Pages[0].Sequence.(*evented.EventPage_Num).Num)
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
