package transport

import "go.uber.org/zap"

type Holder struct {
	Log			*zap.SugaredLogger
	transports  []Transport
	projections []SyncProjection
	sagas       []SyncSaga
}

func (th *Holder) Add(i interface{}) {
	switch i.(type) {
	case Transport:
		th.transports = append(th.transports, i.(Transport))
	default:
		th.Log.Infow("Attempted to add non-transport type to transport Holder.  This may be a synchronous-only transport, and may be OK.")
	}

	switch i.(type) {
	case SyncProjection:
		th.projections = append(th.projections, i.(SyncProjection))
	case SyncSaga:
		th.sagas = append(th.sagas, i.(SyncSaga))
	default:
		th.Log.Infow("Attempted to add non-synchronous type to transport Holder.", "transport", i)
	}
}

func (th *Holder) GetTransports() []Transport {
	return th.transports
}

func (th *Holder) GetProjections() []SyncProjection {
	return th.projections
}

func (th *Holder) GetSaga() []SyncSaga {
	return th.sagas
}

func NewTransportHolder(log *zap.SugaredLogger) *Holder{
	return &Holder{Log: log}
}
