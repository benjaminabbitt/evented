package mock

import (
	"context"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type AsyncTransport struct {
	mock.Mock
}

func (o *AsyncTransport) Handle(ctx context.Context, evts *core.EventBook) (err error) {
	args := o.Called(ctx, evts)
	return args.Error(0)
}
