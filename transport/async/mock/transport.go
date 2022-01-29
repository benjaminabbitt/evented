package mock

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type AsyncTransport struct {
	mock.Mock
}

func (o *AsyncTransport) Handle(ctx context.Context, evts *evented.EventBook) (err error) {
	args := o.Called(ctx, evts)
	return args.Error(0)
}
