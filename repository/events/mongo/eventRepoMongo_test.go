package mongo

import (
	"context"
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type MongoIntegrationSuite struct {
	suite.Suite
	Mongo       events.EventStorer
	populatedId uuid.UUID
	log         *zap.SugaredLogger
	dait        *dockerTestSuite.DockerAssistedIntegrationTest
}

func (s *MongoIntegrationSuite) SetupSuite() {
	s.log = support.Log()
	ctx := context.Background()

	s.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := s.dait.CreateNewContainer("mongo", []uint16{27017})
	if err != nil {
		s.log.Error(err)
	}

	mongo, err := NewEventRepoMongo(fmt.Sprintf("mongodb://localhost:%d", s.dait.Ports[0].PublicPort), "test", "events", s.log)
	if err != nil {
		s.log.Error(err)
	}
	s.Mongo = mongo

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
	_ = s.Mongo.Add(ctx, id, pages)
}

func (s *MongoIntegrationSuite) TearDownSuite() {
	s.log.Sync()
	err := s.dait.StopContainer()
	if err != nil {
		s.log.Error(err)
	}
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
	_ = s.Mongo.Add(context.Background(), id, []*evented_core.EventPage{page})
	result, _ := s.Mongo.Get(context.Background(), id)
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
	_ = s.Mongo.Add(context.Background(), id, []*evented_core.EventPage{page})

	cp := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Force{Force: true},
		CreatedAt:   ts,
		Event:       nil,
		Synchronous: false,
	}

	_ = s.Mongo.Add(context.Background(), id, []*evented_core.EventPage{cp})
	result, _ := s.Mongo.Get(context.Background(), id)
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
	_ = s.Mongo.Add(context.Background(), id, pages)
	result, _ := s.Mongo.Get(context.Background(), id)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, result[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, result[1].Sequence)
}

func (s *MongoIntegrationSuite) Test_GetTo() {
	results, _ := s.Mongo.GetTo(context.Background(), s.populatedId, 1)
	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func (s *MongoIntegrationSuite) Test_GetFrom() {
	results, _ := s.Mongo.GetFrom(context.Background(), s.populatedId, 3)
	s.EqualValues(&evented_core.EventPage_Num{Num: 3}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 4}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func (s *MongoIntegrationSuite) Test_GetFromTo() {
	results, _ := s.Mongo.GetFromTo(context.Background(), s.populatedId, 1, 2)
	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, results[0].Sequence)
	s.EqualValues(&evented_core.EventPage_Num{Num: 2}, results[1].Sequence)
	s.EqualValues(2, len(results))
}

func TestMongoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MongoIntegrationSuite))
}
