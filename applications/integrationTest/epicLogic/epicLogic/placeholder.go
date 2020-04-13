package epicLogic

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewPlaceholderEpicLogic(log *zap.SugaredLogger) PlaceholderEpicLogic {
	return PlaceholderEpicLogic{
		log: log,
	}
}

type PlaceholderEpicLogic struct {
	evented_eventHandler.EventHandlerServer
	eventDomain string
	log         *zap.SugaredLogger
}

func (o *PlaceholderEpicLogic) Handle(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.EventBook, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(uuid)
	cover := eventedcore.Cover{
		Domain: o.eventDomain,
		Root:   &root,
	}
	eb := &eventedcore.EventBook{
		Cover: &cover,
		Pages: []*eventedcore.EventPage{&eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Force{Force: true},
			CreatedAt:   &timestamp.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: nil,
	}
	return eb, nil
}

func (o *PlaceholderEpicLogic) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpc.NewServer()

	evented_eventHandler.RegisterEventHandlerServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
