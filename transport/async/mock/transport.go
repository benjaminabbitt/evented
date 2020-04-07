package mock

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type AsyncTransport struct {
	mock.Mock
}

func (o AsyncTransport) Handle(evts *evented_core.EventBook) (err error) {
	args := o.Called(evts)
	return args.Error(0)
}
