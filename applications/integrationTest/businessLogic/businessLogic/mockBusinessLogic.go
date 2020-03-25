package businessLogic

import (
	"flag"
	"fmt"
	"github.com/benjaminabbitt/evented"
	evented_business "github.com/benjaminabbitt/evented/proto/business"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
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
	log.WithFields(log.Fields{"contextualCommand": in}).Info("Business Logic Handle")
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

	log.WithFields(log.Fields{"eventBook": eventBook}).Info("Business Logic Handle")

	return eventBook, nil
}

func (c *MockBusinessLogicServer) Listen(port uint16){
	lis := c.createListener(port)
	grpcServer := grpc.NewServer()

	evented_business.RegisterBusinessLogicServer(grpcServer, c)
	err := grpcServer.Serve(lis)
	c.errh.LogIfErr(err, "Failed starting server")
}

func (c *MockBusinessLogicServer) createListener(port uint16) net.Listener {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	c.errh.LogIfErr(err, "Failed to Listen")
	return lis
}