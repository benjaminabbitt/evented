package mock

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SnapshotRepo struct {
	mock.Mock
}

func (o SnapshotRepo) Get(id uuid.UUID) (snap *evented_core.Snapshot, err error) {
	args := o.Called(id)
	return args.Get(0).(*evented_core.Snapshot), args.Error(1)
}

func (o SnapshotRepo) Put(id uuid.UUID, snap *evented_core.Snapshot) (err error) {
	args := o.Called(id, snap)
	return args.Error(0)
}
