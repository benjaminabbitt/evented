package snapshot_memory

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type SSMemoryRepository struct {
	store map[string]*evented_core.Snapshot
}

func (repos *SSMemoryRepository) Put(ctx context.Context, id uuid.UUID, ss *evented_core.Snapshot) error {
	repos.store[id.String()] = ss
	return nil
}

func (repos *SSMemoryRepository) Get(ctx context.Context, id uuid.UUID) (*evented_core.Snapshot, error) {
	return repos.store[id.String()], nil
}

func NewSSMemoryRepository() (repos *SSMemoryRepository) {
	repos = &SSMemoryRepository{}
	repos.store = make(map[string]*evented_core.Snapshot)
	return repos
}
