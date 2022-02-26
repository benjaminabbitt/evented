package client

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client *MockClient) Handle(ctx context.Context, command *evented.ContextualCommand) (events *evented.EventBook, err error) {
	args := client.Called(ctx, command)
	return args.Get(0).(*evented.EventBook), args.Error(1)
}
