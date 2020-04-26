package eventBook

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	memoryRepository "github.com/benjaminabbitt/evented/repository/events/event-memory"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type EventBookRepositorySuite struct {
	suite.Suite
	domain string
	id     string
	repos  RepositoryBasic
	book   *evented_core.EventBook
}

func (s *EventBookRepositorySuite) SetupTest() {
	id, _ := uuid.NewRandom()
	s.id = id.String()
	s.domain = "test"

	cover := &evented_core.Cover{
		Domain: s.domain,
		Root:   id.String(),
	}

	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    0,
			CreatedAt:   time.Now().Format(time.RFC3339),
			Event:       nil,
			Synchronous: false,
		},
	}

	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	s.book = &evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}

	s.repos = RepositoryBasic{
		EventRepo:    memoryRepository.NewMemoryRepository(),
		SnapshotRepo: snapshot_memory.NewSSMemoryRepository(),
		Domain:       "test",
	}

	_ = s.repos.Put(s.book)
}

func (s *EventBookRepositorySuite) Test_Put() {
	id, _ := uuid.NewRandom()

	cover := &evented_core.Cover{
		Domain: "testPut",
		Root:   id.String(),
	}

	pages := []*evented_core.EventPage{
		&evented_core.EventPage{
			Sequence:    0,
			CreatedAt:   time.Now().Format(time.RFC3339),
			Event:       nil,
			Synchronous: false,
		},
	}

	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	book := evented_core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: snapshot,
	}

	ebr := &RepositoryBasic{
		EventRepo:    memoryRepository.NewMemoryRepository(),
		SnapshotRepo: snapshot_memory.NewSSMemoryRepository(),
		Domain:       s.domain,
	}

	err := ebr.Put(book)

	s.Assert().NoError(err)
}

func (s *EventBookRepositorySuite) Test_Get() {
	book, err := s.repos.Get(s.id)
	s.Assert().NoError(err)
	s.Assert().EqualValues(s.book, book)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(EventBookRepositorySuite))
}
