package main

import (
	"context"
	"crypto/rand"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
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

func (s *QueryHandlerSuite) SetupTest() {
	s.log = support.Log()
	defer s.log.Sync()
	s.ctx = context.Background()
	s.repos = &mer.EventRepository{}
	s.sut = &DefaultEventQueryServer{
		eventRepos: s.repos,
		log:        s.log,
	}
}

func sumPages(pages []*evented_core.EventPage) uint32 {
	size := uint32(0)
	for _, page := range pages {
		size += uint32(proto.Size(page))
	}
	return size
}

func (o *QueryHandlerSuite) TestMemoryRead() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := evented_proto.UUIDToProto(uuid)
	cover := &evented_core.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}

	book := &evented_core.EventBook{
		Cover:    cover,
		Pages:    nil,
		Snapshot: nil,
	}

	e1 := make([]byte, 16)
	rand.Read(e1)
	pba := &evented_query.TestByteArray{
		Bytes: e1,
	}
	a1, _ := ptypes.MarshalAny(pba)

	e2 := make([]byte, 32)
	rand.Read(e2)
	pba2 := &evented_query.TestByteArray{
		Bytes: e2,
	}
	a2, _ := ptypes.MarshalAny(pba2)

	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 1},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       a1,
			Synchronous: false,
		},
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 2},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       a2,
			Synchronous: false,
		},
	}

	o.log.Info("base: ", proto.Size(book))
	o.log.Info("pages: ", sumPages(pages))

	book2 := &evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	o.log.Info("addition: ", uint32(proto.Size(book))+sumPages(pages))
	o.log.Info("composite: ", proto.Size(book2))
}

func (o *QueryHandlerSuite) Test_Low_High() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := evented_proto.UUIDToProto(uuid)

	//evtChan := make(chan *evented_core.EventPage)
	query := &evented_query.Query{
		Domain:     "test",
		Root:       &protoUUID,
		LowerBound: 1,
		UpperBound: 2,
	}
	ctx := context.Background()
	o.repos.On("GetFromTo", mock.Anything, mock.AnythingOfType("chan *evented_core.EventPage"), uuid, uint32(1), uint32(2)).Return(nil).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan *evented_core.EventPage)
		page := framework.NewEmptyEventPage(0, false)
		go func() {
			ch <- page
			close(ch)
		}()
		o.log.Info("test")
	}).Once()

	queryResponse := MockGetEventsServer{}
	queryResponse.On("Context").Return(ctx)
	queryResponse.On("Send", mock.Anything).Return(nil).Once()

	_ = o.sut.GetEvents(query, queryResponse)

	o.repos.AssertExpectations(o.T())
	queryResponse.AssertExpectations(o.T())
}

//func (o *QueryHandlerSuite) Test_Low() {
//	uuid, _ := uuid2.NewRandom()
//	protoUUID := evented_proto.UUIDToProto(uuid)
//	cover := &evented_core.Cover{
//		Root:   &protoUUID,
//		Domain: "test",
//	}
//	pages := []*evented_core.EventPage{
//		&evented_core.EventPage{
//			Sequence:    &evented_core.EventPage_Num{Num: 1},
//			CreatedAt:   &timestamp.Timestamp{},
//			Event:       nil,
//			Synchronous: false,
//		},
//	}
//	book := evented_core.EventBook{
//		Cover:    cover,
//		Pages:    pages,
//		Snapshot: nil,
//	}
//	o.eventRepos.On("GetFrom", mock.Anything, uuid, uint32(1)).Return(book, nil)
//	query := &evented_query.Query{
//		Domain:     "test",
//		Root:       &protoUUID,
//		LowerBound: 1,
//	}
//	retBook, _ := o.sut.GetEventBook(o.ctx, query)
//	o.eventRepos.AssertExpectations(o.T())
//	o.Assert().Equal(uint32(1), retBook.Pages[0].Sequence.(*evented_core.EventPage_Num).Num)
//}
//
//func (o *QueryHandlerSuite) Test_NoLimits() {
//	uuid, _ := uuid2.NewRandom()
//	protoUUID := evented_proto.UUIDToProto(uuid)
//	cover := &evented_core.Cover{
//		Root:   &protoUUID,
//		Domain: "test",
//	}
//	pages := []*evented_core.EventPage{
//		&evented_core.EventPage{
//			Sequence:    &evented_core.EventPage_Num{Num: 0},
//			CreatedAt:   &timestamp.Timestamp{},
//			Event:       nil,
//			Synchronous: false,
//		},
//		&evented_core.EventPage{
//			Sequence:    &evented_core.EventPage_Num{Num: 1},
//			CreatedAt:   &timestamp.Timestamp{},
//			Event:       nil,
//			Synchronous: false,
//		},
//	}
//	book := evented_core.EventBook{
//		Cover:    cover,
//		Pages:    pages,
//		Snapshot: nil,
//	}
//	o.eventRepos.On("Get", mock.Anything, uuid).Return(book, nil)
//	query := &evented_query.Query{
//		Domain: "test",
//		Root:   &protoUUID,
//	}
//	retBook, _ := o.sut.GetEventBook(o.ctx, query)
//	o.eventRepos.AssertExpectations(o.T())
//	o.Assert().Equal(uint32(0), retBook.Pages[0].Sequence.(*evented_core.EventPage_Num).Num)
//	o.Assert().Equal(uint32(1), retBook.Pages[len(retBook.Pages)-1].Sequence.(*evented_core.EventPage_Num).Num)
//}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(QueryHandlerSuite))
}
