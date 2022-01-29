package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type MockSagaClient struct {
	mock.Mock
}

func (o *MockSagaClient) HandleSync(ctx context.Context, evts *core.EventBook) (responseEvents *core.SynchronousProcessingResponse, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*core.SynchronousProcessingResponse), args.Error(1)
}
