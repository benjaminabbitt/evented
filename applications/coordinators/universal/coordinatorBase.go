package universal

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
)

type Coordinator struct {
	Processed        *processed.Processed
	EventQueryClient evented_query.EventQueryClient
	Log              *zap.SugaredLogger
}

func (o *Coordinator) RepairSequencing(ctx context.Context, eb *evented_core.EventBook, cb func(*evented_core.EventBook) error) {
	id, err := evented_proto.ProtoToUUID(eb.Cover.Root)
	last, err := o.Processed.LastReceived(ctx, id)
	seq := eb.Pages[0].Sequence.(*evented_core.EventPage_Num).Num
	if err != nil {
		o.Log.Error(err)
	}
	nextEventSeq := last + 1
	if nextEventSeq < seq {
		evtStream, err := o.EventQueryClient.GetEvents(ctx, &evented_query.Query{
			Domain:     eb.Cover.Domain,
			Root:       eb.Cover.Root,
			LowerBound: last,
		})
		if err != nil {
			o.Log.Error(err)
		}
		for {
			event, err := evtStream.Recv()
			if err != nil {
				o.Log.Error(err)
			}
			err = cb(event)
			if err != nil {
				o.Log.Error(err)
			} else {
				o.MarkProcessed(ctx, event)
			}
		}
	}
}

func (o *Coordinator) MarkProcessed(ctx context.Context, event *evented_core.EventBook) {
	id, err := evented_proto.ProtoToUUID(event.Cover.Root)
	for _, page := range event.Pages {
		err = o.Processed.Received(ctx, id, page.Sequence.(*evented_core.EventPage_Num).Num)
		if err != nil {
			o.Log.Error(err)
		}
	}
}
