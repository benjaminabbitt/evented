package mock

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type AsyncTransport struct {
	mock.Mock
}

func (o AsyncTransport) Handle(ctx context.Context, evts *evented_core.EventBook) (err error) {
	args := o.Called(ctx, evts)
	return args.Error(0)
}
