package snapshots

import evented_core "github.com/benjaminabbitt/evented/proto/core"

type SnapshotRepo interface {
	Get(id string)(snap *evented_core.Snapshot, err error)
	Put(id string, snap *evented_core.Snapshot)(err error)
}
