package businessLogic

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPlaceholderBusinessLogicServer(log *zap.SugaredLogger) PlaceholderBusinessLogicServer {
	return PlaceholderBusinessLogicServer{
		log: log,
	}
}

type PlaceholderBusinessLogicServer struct {
	evented.UnimplementedBusinessLogicServer
	log *zap.SugaredLogger
}

func (o PlaceholderBusinessLogicServer) Handle(ctx context.Context, in *evented.ContextualCommand) (*evented.EventBook, error) {
	o.log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*evented.EventPage
	//TODO: harden
	ts := timestamppb.Now()
	for _, commandPage := range in.Command.Pages {
		eventPage := &evented.EventPage{
			Sequence:    &evented.EventPage_Num{Num: commandPage.Sequence},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: true,
		}
		eventPages = append(eventPages, eventPage)
	}

	eventBook := &evented.EventBook{
		Cover:    in.Command.Cover,
		Pages:    eventPages,
		Snapshot: nil,
	}

	o.log.Infow("Business Logic Handle", "eventBook", support.StringifyEventBook(eventBook))

	return eventBook, nil
}
