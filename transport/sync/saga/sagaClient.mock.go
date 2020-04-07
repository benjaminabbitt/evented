package saga

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
)

type MockSagaClient struct{
	mock.Mock
}

func (o MockSagaClient) HandleSync(evts *evented_core.EventBook) (responseEvents *evented_core.EventBook, err error) {
	args := o.Called(evts)
	return args.Get(0).(*evented_core.EventBook), args.Error(1)
}

