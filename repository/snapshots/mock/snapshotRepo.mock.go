package mock

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SnapshotRepo struct {
	mock.Mock
}

func (o SnapshotRepo) Get(ctx context.Context, id uuid.UUID) (snap *core.Snapshot, err error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*core.Snapshot), args.Error(1)
}

func (o SnapshotRepo) Put(ctx context.Context, id uuid.UUID, snap *core.Snapshot) (err error) {
	args := o.Called(ctx, id, snap)
	return args.Error(0)
}
