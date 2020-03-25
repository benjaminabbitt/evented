package businessLogic

import (
	"github.com/benjaminabbitt/evented"
	evented_business "github.com/benjaminabbitt/evented/proto/business"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func NewMockBusinessLogic(log *zap.SugaredLogger, errh *evented.ErrLogger) MockBusinessLogicServer {
	return MockBusinessLogicServer{
		log: log,
		errh: errh,
	}
}

type MockBusinessLogicServer struct {
	evented_business.UnimplementedBusinessLogicServer
	log *zap.SugaredLogger
	errh *evented.ErrLogger
}

func (c *MockBusinessLogicServer) Handle(ctx context.Context, in *eventedcore.ContextualCommand) (*eventedcore.EventBook, error){
	c.log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*eventedcore.EventPage
	for _, commandPage := range in.Command.Pages{
		eventPage := &eventedcore.EventPage{
			Sequence:    commandPage.Sequence,
			CreatedAt:   time.Now().Format(time.RFC3339),
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

	c.log.Infow("Business Logic Handle", "eventBook", eventBook)

	return eventBook, nil
}

func (c *MockBusinessLogicServer) Listen(port uint16){
	lis := support.CreateListener(port, c.errh)
	grpcServer := grpc.NewServer()

	evented_business.RegisterBusinessLogicServer(grpcServer, c)
	err := grpcServer.Serve(lis)
	c.errh.LogIfErr(err, "Failed starting server")
}

