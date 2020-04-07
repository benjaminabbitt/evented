package businessLogic

import (
	"github.com/benjaminabbitt/evented"
	evented_business "github.com/benjaminabbitt/evented/proto/business"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func NewSimpleBusinessLogicServer(log *zap.SugaredLogger, errh *evented.ErrLogger) SimpleBusinessLogicServer {
	return SimpleBusinessLogicServer{
		log:  log,
		errh: errh,
	}
}

type SimpleBusinessLogicServer struct {
	evented_business.UnimplementedBusinessLogicServer
	log  *zap.SugaredLogger
	errh *evented.ErrLogger
}

func (c *SimpleBusinessLogicServer) Handle(ctx context.Context, in *eventedcore.ContextualCommand) (*eventedcore.EventBook, error) {
	c.log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*eventedcore.EventPage
	//TODO: harden
	ts, _ := ptypes.TimestampProto(time.Now())
	for _, commandPage := range in.Command.Pages {
		eventPage := &eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Num{Num: commandPage.Sequence},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: false,
		}
		eventPages = append(eventPages, eventPage)
	}

	eventBook := &eventedcore.EventBook{
		Cover:    in.Command.Cover,
		Pages:    eventPages,
		Snapshot: nil,
	}

	c.log.Infow("Business Logic Handle", "eventBook", support.StringifyEventBook(eventBook))

	return eventBook, nil
}

func (c *SimpleBusinessLogicServer) Listen(port uint16) {
	lis := support.CreateListener(port, c.errh)
	grpcServer := grpc.NewServer()

	evented_business.RegisterBusinessLogicServer(grpcServer, c)
	err := grpcServer.Serve(lis)
	c.errh.LogIfErr(err, "Failed starting server")
}
