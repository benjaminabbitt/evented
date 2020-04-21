package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type MockSagaClient struct {
	mock.Mock
}

func (o *MockSagaClient) HandleSync(ctx context.Context, evts *evented_core.EventBook) (responseEvents *evented_core.SynchronousProcessingResponse, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*evented_core.SynchronousProcessingResponse), args.Error(1)
}
