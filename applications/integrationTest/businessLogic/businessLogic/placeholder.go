package businessLogic

import (
	eventedbusiness "github.com/benjaminabbitt/evented/proto/evented/business"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

func NewPlaceholderBusinessLogicServer(log *zap.SugaredLogger) PlaceholderBusinessLogicServer {
	return PlaceholderBusinessLogicServer{
		log: log,
	}
}

type PlaceholderBusinessLogicServer struct {
	eventedbusiness.UnimplementedBusinessLogicServer
	log *zap.SugaredLogger
}

func (o PlaceholderBusinessLogicServer) Handle(ctx context.Context, in *eventedcore.ContextualCommand) (*eventedcore.EventBook, error) {
	o.log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*eventedcore.EventPage
	//TODO: harden
	ts, _ := ptypes.TimestampProto(time.Now())
	for _, commandPage := range in.Command.Pages {
		eventPage := &eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Num{Num: commandPage.Sequence},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: true,
		}
		eventPages = append(eventPages, eventPage)
	}

	eventBook := &eventedcore.EventBook{
		Cover:    in.Command.Cover,
		Pages:    eventPages,
		Snapshot: nil,
	}

	o.log.Infow("Business Logic Handle", "eventBook", support.StringifyEventBook(eventBook))

	return eventBook, nil
}

