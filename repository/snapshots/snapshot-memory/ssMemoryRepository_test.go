package snapshot_memory

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SSMemoryRepoSuite struct {
	suite.Suite
	mr *SSMemoryRepository
	id string
}

func (s *SSMemoryRepoSuite) SetupTest() {
	s.id = "a"
	s.mr = NewSSMemoryRepository()
}

func (s *SSMemoryRepoSuite) TestAdd() {
	snapshot := evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	err := s.mr.Put(s.id, &snapshot)
	s.Assert().NoError(err)
}

func (s *SSMemoryRepoSuite) TestGet() {
	snapshot := evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}
	err := s.mr.Put(s.id, &snapshot)
	s.Assert().NoError(err)

	retrieved, err := s.mr.Get(s.id)
	s.Assert().Equal(&snapshot, retrieved)
	s.Assert().NoError(err)
}

func Test(t *testing.T) {
	suite.Run(t, new(SSMemoryRepoSuite))
}
