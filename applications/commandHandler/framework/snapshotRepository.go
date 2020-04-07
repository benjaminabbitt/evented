package framework

import evented_core "github.com/benjaminabbitt/evented/proto/core"

type SnapshotRepository interface {
	Add(evt evented_core.Snapshot) (err error)
	Get(id string) (ent evented_core.Snapshot, err error)
}
