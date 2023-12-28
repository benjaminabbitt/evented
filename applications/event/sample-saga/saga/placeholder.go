package saga

import (
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPlaceholderSagaLogic(log *zap.SugaredLogger, tracer *opentracing.Tracer) PlaceholderSagaLogic {
	return PlaceholderSagaLogic{
		log:    log,
		tracer: tracer,
	}
}

type PlaceholderSagaLogic struct {
	evented2.UnimplementedSagaServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *PlaceholderSagaLogic) Handle(ctx context.Context, in *evented2.EventBook) (*emptypb.Empty, error) {
	_, err := o.HandleSync(ctx, in)
	return &emptypb.Empty{}, err
}

func (o *PlaceholderSagaLogic) HandleSync(ctx context.Context, in *evented2.EventBook) (resp *evented2.SynchronousProcessingResponse, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(id)
	cover := evented2.Cover{
		Domain: o.eventDomain,
		Root:   &root,
	}
	eb := &evented2.EventBook{
		Cover: &cover,
		Pages: []*evented2.EventPage{{
			Sequence:    &evented2.EventPage_Force{Force: true},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: nil,
	}

	resp = &evented2.SynchronousProcessingResponse{
		Books:       []*evented2.EventBook{eb},
		Projections: nil,
	}

	return resp, nil
}

func (o *PlaceholderSagaLogic) Listen(port uint) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	evented2.RegisterSagaServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
