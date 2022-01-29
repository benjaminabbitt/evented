package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type MockProjectorClient struct {
	mock.Mock
}

func (o MockProjectorClient) HandleSync(ctx context.Context, evts *core.EventBook) (projection *core.Projection, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*core.Projection), args.Error(1)
}
