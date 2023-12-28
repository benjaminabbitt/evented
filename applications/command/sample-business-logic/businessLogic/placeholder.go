package businessLogic

import (
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
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
	evented2.UnimplementedBusinessLogicServer
	log *zap.SugaredLogger
}

func (o PlaceholderBusinessLogicServer) Handle(ctx context.Context, in *evented2.ContextualCommand) (*evented2.EventBook, error) {
	o.log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*evented2.EventPage
	//TODO: harden
	ts := timestamppb.Now()
	for _, commandPage := range in.Command.Pages {
		eventPage := &evented2.EventPage{
			Sequence:    &evented2.EventPage_Num{Num: commandPage.Sequence},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: true,
		}
		eventPages = append(eventPages, eventPage)
	}

	eventBook := &evented2.EventBook{
		Cover:    in.Command.Cover,
		Pages:    eventPages,
		Snapshot: nil,
	}

	o.log.Infow("Business Logic Handle", "eventBook", support.StringifyEventBook(eventBook))

	return eventBook, nil
}
