package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/stretchr/testify/mock"
)

type MockProjectorClient struct {
	mock.Mock
}

func (o *MockProjectorClient) HandleSync(ctx context.Context, evts *evented.EventBook) (projection *evented.Projection, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*evented.Projection), args.Error(1)
}
