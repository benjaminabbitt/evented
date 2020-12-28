package memory
//
//import (
//	"context"
//	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
//	"github.com/benjaminabbitt/evented/repository/events"
//	"github.com/benjaminabbitt/evented/support"
//	"github.com/golang/protobuf/ptypes"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/suite"
//	"go.uber.org/zap"
//	"testing"
//	"time"
//)
//
//type MemorySuite struct {
//	suite.Suite
//	Repo        events.EventStorer
//	populatedId uuid.UUID
//	log         *zap.SugaredLogger
//}
//
//func (s *MemorySuite) SetupSuite() {
//	s.log = support.Log()
//	ctx := context.Background()
//
//	repo, err := NewEventRepoMemory(s.log)
//	if err != nil {
//		s.log.Error(err)
//	}
//	s.Repo = repo
//
//	ts, _ := ptypes.TimestampProto(time.Now())
//	id, _ := uuid.NewRandom()
//	s.populatedId = id
//
//	var pages []*evented_core.EventPage
//
//	pages = append(pages, &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Num{Num: 0},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	})
//	for i := 0; i < 4; i++ {
//		pages = append(pages, &evented_core.EventPage{
//			Sequence:    &evented_core.EventPage_Force{Force: true},
//			CreatedAt:   ts,
//			Event:       nil,
//			Synchronous: false,
//		})
//	}
//	_ = s.Repo.Add(ctx, id, pages)
//}
//
//func (s *MemorySuite) TearDownSuite() {
//	s.log.Sync()
//}
//
//func (s *MemorySuite) Test_Insert_Sequence() {
//	ts, _ := ptypes.TimestampProto(time.Now())
//	id, _ := uuid.NewRandom()
//	page := &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Num{Num: 0},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	}
//	_ = s.Repo.Add(context.Background(), id, []*evented_core.EventPage{page})
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.Get(context.Background(), ch, id)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, (<-ch).Sequence)
//}
//
//func (s *MemorySuite) Test_Insert_Force_Preexisting_Sequence() {
//	ts, _ := ptypes.TimestampProto(time.Now())
//	id, _ := uuid.NewRandom()
//	page := &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Force{Force: true},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	}
//	_ = s.Repo.Add(context.Background(), id, []*evented_core.EventPage{page})
//
//	cp := &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Force{Force: true},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	}
//
//	_ = s.Repo.Add(context.Background(), id, []*evented_core.EventPage{cp})
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.Get(context.Background(), ch, id)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, (<-ch).Sequence)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, (<-ch).Sequence)
//}
//
//func (s *MemorySuite) Test_Force_With_Numbered_In_Same_Book() {
//	ts, _ := ptypes.TimestampProto(time.Now())
//	id, _ := uuid.NewRandom()
//
//	var pages []*evented_core.EventPage
//
//	pages = append(pages, &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Num{Num: 0},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	})
//	pages = append(pages, &evented_core.EventPage{
//		Sequence:    &evented_core.EventPage_Force{Force: true},
//		CreatedAt:   ts,
//		Event:       nil,
//		Synchronous: false,
//	})
//	_ = s.Repo.Add(context.Background(), id, pages)
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.Get(context.Background(), ch, id)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, (<-ch).Sequence)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, (<-ch).Sequence)
//}
//
//func (s *MemorySuite) Test_GetTo() {
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.GetTo(context.Background(), ch, s.populatedId, 3)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 0}, (<-ch).Sequence)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, (<-ch).Sequence)
//	s.Assert().Empty(ch)
//}
//
//func (s *MemorySuite) Test_GetFrom() {
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.GetFrom(context.Background(), ch, s.populatedId, 3)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 3}, (<-ch).Sequence)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 4}, (<-ch).Sequence)
//	s.Assert().Empty(ch)
//}
//
//func (s *MemorySuite) Test_GetFromTo() {
//	ch := make(chan *evented_core.EventPage)
//	_ = s.Repo.GetFromTo(context.Background(), ch, s.populatedId, 1, 3)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 1}, (<-ch).Sequence)
//	s.EqualValues(&evented_core.EventPage_Num{Num: 2}, (<-ch).Sequence)
//	s.Assert().Empty(ch)
//}
//
//func TestMongoIntegrationSuite(t *testing.T) {
//	if !testing.Short() {
//		suite.Run(t, new(MemorySuite))
//	}
//}
