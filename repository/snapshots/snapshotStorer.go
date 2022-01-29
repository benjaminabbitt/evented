package snapshots

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
)

type SnapshotStorer interface {
	Get(ctx context.Context, id uuid.UUID) (snap *core.Snapshot, err error)
	Put(ctx context.Context, id uuid.UUID, snap *core.Snapshot) (err error)
}
