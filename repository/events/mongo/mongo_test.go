package mongo

import (
	"github.com/benjaminabbitt/evented"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

type MongoIntegrationSuite struct {
	suite.Suite
	Mongo       *Mongo
	populatedId uuid.UUID
}

func (s *MongoIntegrationSuite) SetupTest() {
	log, errh = support.Log()
	defer log.Sync()
	s.Mongo = NewMongoClient("mongodb://localhost:27017", "test", "events", log, errh)

	ts, _ := ptypes.TimestampProto(time.Now())
	id, _ := uuid.NewRandom()
	s.populatedId = id

	var pages []*evented_core.EventPage

	pages = append(pages, &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	})
	for i := 0; i < 4; i++ {
		pages = append(pages, &evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Force{Force: true},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: false,
		})
	}
	_ = s.Mongo.Add(id, pages)
}

func (s *MongoIntegrationSuite) Test_Insert_Sequence() {
	ts, _ := ptypes.TimestampProto(time.Now())
	id, _ := uuid.NewRandom()
	page := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	}
	_ = s.Mongo.Add(id, []*evented_core.EventPage{page})
	result, _ := s.Mongo.Get(id)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, result[0].Sequence)
}

func (s *MongoIntegrationSuite) Test_Insert_Force_Preexisting_Sequence() {
	ts, _ := ptypes.TimestampProto(time.Now())
	id, _ := uuid.NewRandom()
	page := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Force{Force: true},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	}
	_ = s.Mongo.Add(id, []*evented_core.EventPage{page})
	_ = s.Mongo.Add(id, []*evented_core.EventPage{page})
	result, _ := s.Mongo.Get(id)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, result[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, result[1].Sequence)
}

func (s *MongoIntegrationSuite) Test_Force_With_Numbered_In_Same_Book() {
	ts, _ := ptypes.TimestampProto(time.Now())
	id, _ := uuid.NewRandom()

	var pages []*evented_core.EventPage

	pages = append(pages, &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	})
	pages = append(pages, &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Force{Force: true},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	})
	_ = s.Mongo.Add(id, pages)
	result, _ := s.Mongo.Get(id)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, result[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, result[1].Sequence)
}

func (s *MongoIntegrationSuite) Test_Get_Next_Sequence_No_Preexisting() {
	id, _ := uuid.NewRandom()
	next, _ := s.Mongo.GetNextSequence(id)
	s.EqualValues(0, next.Sequence)
}

func (s *MongoIntegrationSuite) Test_Get_Next_Sequence_Preexisting() {
	ts, _ := ptypes.TimestampProto(time.Now())
	id, _ := uuid.NewRandom()
	page := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	}
	_ = s.Mongo.Add(id, []*evented_core.EventPage{page})
	next, _ := s.Mongo.GetNextSequence(id)
	s.EqualValues(1, next.Sequence)
}

func (s *MongoIntegrationSuite) Test_GetTo() {
	results, _ := s.Mongo.GetTo(s.populatedId, 1)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func (s *MongoIntegrationSuite) Test_GetFrom() {
	results, _ := s.Mongo.GetFrom(s.populatedId, 3)
	s.EqualValues(&evented_core.EventPage_Num{Num: 3}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 4}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func (s *MongoIntegrationSuite) Test_GetFromTo() {
	results, _ := s.Mongo.GetFromTo(s.populatedId, 1, 2)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 2}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func TestMongoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MongoIntegrationSuite))
}
