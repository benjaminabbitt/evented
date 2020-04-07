package snapshots

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

type SnapshotRepo interface {
	Get(id uuid.UUID) (snap *evented_core.Snapshot, err error)
	Put(id uuid.UUID, snap *evented_core.Snapshot) (err error)
}
