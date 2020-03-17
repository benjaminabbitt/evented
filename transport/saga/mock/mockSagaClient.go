package mock

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport"
	"github.com/stretchr/testify/mock"
)

type SagaClient struct{
	mock.Mock
}

func (client SagaClient) SendSync(evts *evented_core.EventBook)(responseEvents *evented_core.EventBook, err error){
	client.Called(evts)

	eventBook := &evented_core.EventBook{
		Cover:    evts.Cover,
		Pages:    nil,
		Snapshot: nil,
	}

	return eventBook, nil
}

func (client SagaClient) Send (evts *evented_core.EventBook)(err error){
	client.Called(evts)
	return nil
}

func NewSagaClient() transport.Saga {
	return &SagaClient{}
}