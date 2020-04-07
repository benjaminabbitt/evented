package event_memory

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"time"
)
import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MemoryRepoSuite struct {
	suite.Suite
	mr MemoryRepository
	id string
}

func (s *MemoryRepoSuite) SetupTest() {
	s.id = "a"
	s.mr = NewMemoryRepository()
	se0 := evented_core.EventPage{
		Sequence:    0,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Event:       nil,
		Synchronous: false,
	}
	se1 := evented_core.EventPage{
		Sequence:    1,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Event:       nil,
		Synchronous: false,
	}
	se2 := evented_core.EventPage{
		Sequence:    2,
		CreatedAt:   time.Now().Format(time.RFC3339),
		Event:       nil,
		Synchronous: false,
	}
	s.mr.Add(s.id, []*evented_core.EventPage{&se0, &se1, &se2})
}

func (s *MemoryRepoSuite) TestGetEventsTo() {
	events, _ := s.mr.GetTo(s.id, 1)
	if len(events) != 2 {
		s.T().Fail()
	}
	s.Assert().Equal(uint32(0), events[0].Sequence)
	s.Assert().Equal(uint32(1), events[1].Sequence)
}

func (s *MemoryRepoSuite) TestGetEventsFrom() {
	events, _ := s.mr.GetFrom(s.id, 1)
	if len(events) != 2 {
		s.T().Fail()
	}
	s.Assert().Equal(uint32(1), events[0].Sequence)
	s.Assert().Equal(uint32(2), events[1].Sequence)
}

func (s *MemoryRepoSuite) TestGetEventsFromTo() {
	events, _ := s.mr.GetFromTo(s.id, 1, 1)
	if len(events) != 1 {
		s.T().Fail()
	}
	s.Assert().Equal(uint32(1), events[0].Sequence)
}

func Test(t *testing.T) {
	suite.Run(t, new(MemoryRepoSuite))
}
