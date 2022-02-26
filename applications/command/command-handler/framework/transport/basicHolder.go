package transport

import (
	"errors"
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"go.uber.org/zap"
	"reflect"
)

type BasicHolder struct {
	Log        *zap.SugaredLogger
	transports []chan *evented.EventBook
	projectors []evented.ProjectorClient
	sagas      []saga.SyncSagaTransporter
}

func (th *BasicHolder) Add(i interface{}) error {
	switch i.(type) {
	case chan *evented.EventBook:
		th.transports = append(th.transports, i.(chan *evented.EventBook))
	case evented.ProjectorClient:
		th.projectors = append(th.projectors, i.(evented.ProjectorClient))
	case saga.SyncSagaTransporter:
		th.sagas = append(th.sagas, i.(saga.SyncSagaTransporter))
	default:
		return errors.New(fmt.Sprintf("Attempted to add unrecognized type %s to transport BasicHolder.", reflect.TypeOf(i).Name()))
	}
	return nil
}

func (th *BasicHolder) GetTransports() []chan *evented.EventBook {
	return th.transports
}

func (th *BasicHolder) GetProjectors() []evented.ProjectorClient {
	return th.projectors
}

func (th *BasicHolder) GetSaga() []saga.SyncSagaTransporter {
	return th.sagas
}

func NewTransportHolder(log *zap.SugaredLogger) *BasicHolder {
	return &BasicHolder{Log: log}
}
