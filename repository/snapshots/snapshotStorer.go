package snapshots

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/google/uuid"
)

type SnapshotStorer interface {
	Get(ctx context.Context, id uuid.UUID) (snap *evented.Snapshot, err error)
	Put(ctx context.Context, id uuid.UUID, snap *evented.Snapshot) (err error)
}
