package mock

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SnapshotRepo struct {
	mock.Mock
}

func (o *SnapshotRepo) Get(ctx context.Context, id uuid.UUID) (snap *evented.Snapshot, err error) {
	args := o.Called(ctx, id)
	return args.Get(0).(*evented.Snapshot), args.Error(1)
}

func (o *SnapshotRepo) Put(ctx context.Context, id uuid.UUID, snap *evented.Snapshot) (err error) {
	args := o.Called(ctx, id, snap)
	return args.Error(0)
}
