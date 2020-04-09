package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type MockProjectorClient struct {
	mock.Mock
}

func (o MockProjectorClient) HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*evented_core.Projection), args.Error(1)
}