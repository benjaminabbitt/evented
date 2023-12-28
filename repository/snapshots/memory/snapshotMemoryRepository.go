package memory

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SnapshotMemoryRepo struct {
	log   *zap.SugaredLogger
	store map[uuid.UUID]*evented.Snapshot
}

func (o SnapshotMemoryRepo) Get(ctx context.Context, id uuid.UUID) (snap *evented.Snapshot, err error) {
	return o.store[id], nil
}

func (o SnapshotMemoryRepo) Put(ctx context.Context, id uuid.UUID, snap *evented.Snapshot) (err error) {
	o.store[id] = snap
	return nil
}

func NewSnapshotRepoMemory(log *zap.SugaredLogger) (snapshots.SnapshotStorer, error) {
	return SnapshotMemoryRepo{
		store: make(map[uuid.UUID]*evented.Snapshot),
		log:   log,
	}, nil
}
