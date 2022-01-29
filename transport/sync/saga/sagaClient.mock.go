package saga

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/stretchr/testify/mock"
)

type MockSagaClient struct {
	mock.Mock
}

func (o *MockSagaClient) HandleSync(ctx context.Context, evts *evented.EventBook) (responseEvents *evented.SynchronousProcessingResponse, err error) {
	args := o.Called(ctx, evts)
	return args.Get(0).(*evented.SynchronousProcessingResponse), args.Error(1)
}
