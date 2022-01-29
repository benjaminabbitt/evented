package memory

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SnapshotMemoryRepo struct {
	log   *zap.SugaredLogger
	store map[uuid.UUID]*core.Snapshot
}

func (o SnapshotMemoryRepo) Get(ctx context.Context, id uuid.UUID) (snap *core.Snapshot, err error) {
	return o.store[id], nil
}

func (o SnapshotMemoryRepo) Put(ctx context.Context, id uuid.UUID, snap *core.Snapshot) (err error) {
	o.store[id] = snap
	return nil
}

func NewSnapshotRepoMemory(log *zap.SugaredLogger) (snapshots.SnapshotStorer, error) {
	return SnapshotMemoryRepo{
		store: make(map[uuid.UUID]*core.Snapshot),
		log:   log,
	}, nil
}
