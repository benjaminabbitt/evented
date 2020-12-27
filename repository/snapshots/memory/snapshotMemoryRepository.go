package memory

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SnapshotMongoRepo struct {
	log   *zap.SugaredLogger
	store map[uuid.UUID]*evented_core.Snapshot
}

func (o *SnapshotMongoRepo) Get(ctx context.Context, id uuid.UUID) (snap *evented_core.Snapshot, err error) {
	return o.store[id], nil
}

func (o *SnapshotMongoRepo) Put(ctx context.Context, id uuid.UUID, snap *evented_core.Snapshot) (err error) {
	o.store[id] = snap
	return nil
}
