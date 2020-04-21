package projector

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	evented_saga "github.com/benjaminabbitt/evented/proto/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewPlaceholderProjectorLogic(log *zap.SugaredLogger) PlaceholderSagaLogic {
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
	_, err := o.HandleSync(ctx, in)
	return &empty.Empty{}, err
}

func (o *PlaceholderSagaLogic) HandleSync(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.Projection, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(id)
	cover := &eventedcore.Cover{
		Domain: o.eventDomain,
		Root:   &root,
	}
	eb := &eventedcore.Projection{
		Cover:      cover,
		Projector:  "test",
		Sequence:   in.Pages[len(in.Pages)].Sequence.(*eventedcore.EventPage_Num).Num,
		Projection: nil,
	}
	return eb, nil
}

func (o *PlaceholderSagaLogic) Listen(port uint16) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar())

	evented_projector.RegisterProjectorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
