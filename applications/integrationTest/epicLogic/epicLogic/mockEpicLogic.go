package epicLogic

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_eventHandler "github.com/benjaminabbitt/evented/proto/eventHandler"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

func NewMockEpicLogic(log *zap.SugaredLogger, errh *evented.ErrLogger) MockEpicLogicServer {
	return MockEpicLogicServer{
		log: log,
		errh: errh,
	}
}

type MockEpicLogicServer struct {
	evented_eventHandler.EventHandlerServer
	log *zap.SugaredLogger
	errh *evented.ErrLogger
}

func (c *MockEpicLogicServer) Handle(ctx context.Context, in *eventedcore.EventBook)(*eventedcore.EventBook, error){
	uuid, err := uuid.NewRandom()
	c.log.Warn(err)
	root := evented_proto.UUIDToProto(uuid)
	cover := eventedcore.Cover{
		Domain: "domain2",
		Root:   &root,
	}
	eb := &eventedcore.EventBook{
		Cover:    &cover,
		Pages:    nil,
		Snapshot: nil,
	}
	return eb, nil
}

func (c *MockEpicLogicServer) Listen(port uint16){
	lis := c.createListener(port)
	grpcServer := grpc.NewServer()

	evented_eventHandler.RegisterEventHandlerServer(grpcServer, c)
	err := grpcServer.Serve(lis)
	c.errh.LogIfErr(err, "Failed starting server")
}

func (c *MockEpicLogicServer) createListener(port uint16) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	c.errh.LogIfErr(err, "Failed to Listen")
	return lis
}