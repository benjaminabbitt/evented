package universal

import (
	"context"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
)

type Coordinator struct {
	Processed        *processed.Processed
	EventQueryClient query.EventQueryClient
	Log              *zap.SugaredLogger
}

func (o *Coordinator) RepairSequencing(ctx context.Context, eb *core.EventBook, cb func(*core.EventBook) error) {
	id, err := eventedproto.ProtoToUUID(eb.Cover.Root)
	last, err := o.Processed.LastReceived(ctx, id)
	seq := eb.Pages[0].Sequence.(*core.EventPage_Num).Num
	if err != nil {
		o.Log.Error(err)
	}
	nextEventSeq := last + 1
	if nextEventSeq < seq {
		evtStream, err := o.EventQueryClient.GetEvents(ctx, &query.Query{
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

func (o *Coordinator) MarkProcessed(ctx context.Context, event *core.EventBook) {
	id, err := eventedproto.ProtoToUUID(event.Cover.Root)
	for _, page := range event.Pages {
		err = o.Processed.Received(ctx, id, page.Sequence.(*core.EventPage_Num).Num)
		if err != nil {
			o.Log.Error(err)
		}
	}
}
