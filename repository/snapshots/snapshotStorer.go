package snapshots

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type SnapshotStorer interface {
	Get(ctx context.Context, id uuid.UUID) (snap *evented_core.Snapshot, err error)
	Put(ctx context.Context, id uuid.UUID, snap *evented_core.Snapshot) (err error)
}
