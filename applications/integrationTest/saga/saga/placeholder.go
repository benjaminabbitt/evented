package saga

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcZap"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewPlaceholderEpicLogic(log *zap.SugaredLogger) PlaceholderSagaLogic {
	return PlaceholderSagaLogic{
		log: log,
	}
}

type PlaceholderSagaLogic struct {
	evented_saga.UnimplementedSagaServer
	eventDomain string
	log         *zap.SugaredLogger
}

func (o *PlaceholderSagaLogic) Handle(ctx context.Context, in *eventedcore.EventBook) (*empty.Empty, error) {
	_, err := o.Handle(ctx, in)
	return &empty.Empty{}, err
}

func (o *PlaceholderSagaLogic) HandleSync(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.EventBook, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(id)
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

func (o *PlaceholderSagaLogic) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcZap.GenerateConfiguredServer(o.log.Desugar())

	evented_saga.RegisterSagaServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
