package mock

import (
	"github.com/benjaminabbitt/evented/proto/core"
	"go.uber.org/zap"
)

type SagaClient struct{
	log *zap.SugaredLogger
}

func (client SagaClient) HandleSync(evts *evented_core.EventBook)(responseEvents *evented_core.EventBook, err error){
	eventBook := &evented_core.EventBook{
		Cover:    evts.Cover,
		Pages:    nil,
		Snapshot: nil,
	}

	return eventBook, nil
}

func (client SagaClient) Handle (evts *evented_core.EventBook)(err error){
	return nil
}

func NewSagaClient(log *zap.SugaredLogger) *SagaClient {
	return &SagaClient{log}
}