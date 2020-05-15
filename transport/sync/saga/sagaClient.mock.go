package saga

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type MockSagaClient struct {
	mock.Mock
}

func (o *MockSagaClient) HandleSync(ctx context.Context, evts *eventedcore.EventBook) (responseEvents *eventedcore.SynchronousProcessingResponse, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*eventedcore.SynchronousProcessingResponse), args.Error(1)
}
