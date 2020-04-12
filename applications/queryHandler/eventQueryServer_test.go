package main

import (
	"context"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
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
	errh  *evented.ErrLogger
	ctx   context.Context
	repos *eventBook.MockEventBookRepository
	sut   EventQueryServer
}

func (s *QueryHandlerSuite) SetupTest() {
	s.log, s.errh = support.Log()
	defer s.log.Sync()
	s.ctx = context.Background()
	s.repos = &eventBook.MockEventBookRepository{}
	s.sut = EventQueryServer{
		repos: s.repos,
		log:   s.log,
		errh:  s.errh,
	}
}

func (o *QueryHandlerSuite) Test_Low_High() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := evented_proto.UUIDToProto(&uuid)
	cover := &evented_core.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}
	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 1},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 2},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
	}
	book := evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	o.repos.On("GetFromTo", mock.Anything, uuid, uint32(1), uint32(2)).Return(book, nil)
	query := &evented_query.Query{
		Domain:     "test",
		Root:       &protoUUID,
		LowerBound: 1,
		UpperBound: 2,
	}
	retBook, _ := o.sut.GetEventBook(o.ctx, query)
	o.repos.AssertExpectations(o.T())
	o.Assert().Equal(uint32(1), retBook.Pages[0].Sequence.(*evented_core.EventPage_Num).Num)
	o.Assert().Equal(uint32(2), retBook.Pages[len(retBook.Pages)-1].Sequence.(*evented_core.EventPage_Num).Num)
}
func (o *QueryHandlerSuite) Test_Low() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := evented_proto.UUIDToProto(&uuid)
	cover := &evented_core.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}
	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 1},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
	}
	book := evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	o.repos.On("GetFrom", mock.Anything, uuid, uint32(1)).Return(book, nil)
	query := &evented_query.Query{
		Domain:     "test",
		Root:       &protoUUID,
		LowerBound: 1,
	}
	retBook, _ := o.sut.GetEventBook(o.ctx, query)
	o.repos.AssertExpectations(o.T())
	o.Assert().Equal(uint32(1), retBook.Pages[0].Sequence.(*evented_core.EventPage_Num).Num)
}

func (o *QueryHandlerSuite) Test_NoLimits() {
	uuid, _ := uuid2.NewRandom()
	protoUUID := evented_proto.UUIDToProto(&uuid)
	cover := &evented_core.Cover{
		Root:   &protoUUID,
		Domain: "test",
	}
	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 0},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
		&evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: 1},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		},
	}
	book := evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	o.repos.On("Get", mock.Anything, uuid).Return(book, nil)
	query := &evented_query.Query{
		Domain: "test",
		Root:   &protoUUID,
	}
	retBook, _ := o.sut.GetEventBook(o.ctx, query)
	o.repos.AssertExpectations(o.T())
	o.Assert().Equal(uint32(0), retBook.Pages[0].Sequence.(*evented_core.EventPage_Num).Num)
	o.Assert().Equal(uint32(1), retBook.Pages[len(retBook.Pages)-1].Sequence.(*evented_core.EventPage_Num).Num)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(QueryHandlerSuite))
}
