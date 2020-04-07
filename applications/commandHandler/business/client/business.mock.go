package client

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client MockClient) Handle(command *evented_core.ContextualCommand) (events *evented_core.EventBook, err error) {
	args := client.Called(command)
	return args.Get(0).(*evented_core.EventBook), args.Error(1)
}