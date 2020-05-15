package projector

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type MockProjectorClient struct {
	mock.Mock
}

func (o MockProjectorClient) HandleSync(ctx context.Context, evts *eventedcore.EventBook) (projection *eventedcore.Projection, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*eventedcore.Projection), args.Error(1)
}
