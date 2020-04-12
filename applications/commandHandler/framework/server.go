package framework

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func NewServer(eventBookRepository eventBook.Repository, transports transport.TransportHolder, businessClient client.BusinessClient, log *zap.SugaredLogger, errh *evented.ErrLogger) Server {
	return Server{
		errh:                errh,
		log:                 log,
		eventBookRepository: eventBookRepository,
		transports:          transports,
		businessClient:      businessClient,
	}
}

func (server *Server) Listen(port uint16) {
	server.log.Infow("Opening port", "port", port)
	lis := server.createListener(port)
	server.log.Infow("Port opened", "port", port)
	server.log.Infow("Creating GRPC Server")
	grpcServer := grpc.NewServer()
	server.log.Infow("Registering Command Handler with GRPC")
	evented_core.RegisterCommandHandlerServer(grpcServer, server)
	server.log.Infow("Handler registered.")
	server.log.Infow("Serving...")
	err := grpcServer.Serve(lis)
	server.errh.LogIfErr(err, "Failed starting server")
}

func (server *Server) createListener(port uint16) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	server.errh.LogIfErr(err, "Failed to Listen")
	return lis
}

type Server struct {
	evented_core.UnimplementedCommandHandlerServer
	errh                *evented.ErrLogger
	log                 *zap.SugaredLogger
	eventBookRepository eventBook.Repository
	transports          transport.TransportHolder
	businessClient      client.BusinessClient
}

func (server Server) Handle(ctx context.Context, in *evented_core.CommandBook) (result *evented_core.CommandHandlerResponse, err error) {
	uuid, err := evented_proto.ProtoToUUID(in.Cover.Root)
	if err != nil {
		return nil, err
	}
	priorState, err := server.eventBookRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	contextualCommand := &evented_core.ContextualCommand{
		Events:  priorState,
		Command: in,
	}

	businessResponse, err := server.businessClient.Handle(ctx, contextualCommand)
	server.errh.LogIfErr(err, "")
	response, err := server.handleEventBook(ctx, businessResponse)
	return &response, err
}

func (server Server) handleEventBook(ctx context.Context, eb *evented_core.EventBook) (result evented_core.CommandHandlerResponse, err error) {
	err = server.eventBookRepository.Put(ctx, eb)
	server.errh.LogIfErr(err, "")

	sync, _ := server.extractSynchronous(*eb)
	var eventBooks []*evented_core.EventBook
	var projections []*evented_core.Projection

	for _, syncProjector := range server.transports.GetProjections() {
		response, err := syncProjector.HandleSync(ctx, &sync)
		server.errh.LogIfErr(err, "")
		projections = append(projections, response)
	}

	for _, syncSaga := range server.transports.GetSaga() {
		response, err := syncSaga.HandleSync(ctx, &sync)
		server.errh.LogIfErr(err, "")
		eventBooks = append(eventBooks, response)
	}

	for _, t := range server.transports.GetTransports() {
		err := t.Handle(ctx, eb)
		server.errh.LogIfErr(err, "")
	}

	return result, nil
}

func (server Server) extractSynchronous(originalBook evented_core.EventBook) (synchronous evented_core.EventBook, async evented_core.EventBook) {
	var lastIdx uint32
	for idx, event := range originalBook.Pages {
		if event.Synchronous {
			lastIdx = uint32(idx) + 1
		}
	}
	synchronous.Pages = originalBook.Pages[:lastIdx]
	synchronous.Cover = originalBook.Cover
	synchronous.Snapshot = originalBook.Snapshot

	async.Pages = originalBook.Pages[lastIdx:]
	async.Cover = originalBook.Cover
	async.Snapshot = nil

	return synchronous, async
}

func (server Server) Record(ctx context.Context, in *evented_core.EventBook) (response *evented_core.CommandHandlerResponse, err error) {
	r, err := server.handleEventBook(ctx, in)
	return &r, err
}
