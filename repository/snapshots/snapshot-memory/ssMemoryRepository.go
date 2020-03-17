package snapshot_memory

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
)

type SSMemoryRepository struct {
	store map[string]*evented_core.Snapshot
}

func (repos *SSMemoryRepository) Put(id string, ss *evented_core.Snapshot) error {
	repos.store[id] = ss
	return nil
}

func (repos *SSMemoryRepository) Get(id string) (*evented_core.Snapshot, error) {
	return repos.store[id], nil
}

func NewSSMemoryRepository() (repos *SSMemoryRepository) {
	repos = &SSMemoryRepository{}
	repos.store = make(map[string]*evented_core.Snapshot)
	return repos
}
