package mock

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/stretchr/testify/mock"
)

type AsyncTransport struct {
	mock.Mock
}

func (o *AsyncTransport) Handle(ctx context.Context, evts *evented.EventBook) (err error) {
	args := o.Called(ctx, evts)
	return args.Error(0)
}
